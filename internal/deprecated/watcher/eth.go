package watcher

import (
	"context"
	"fmt"
	app "go-blocker/internal/application/payment"
	"go-blocker/internal/deprecated/storage"
	blockchain "go-blocker/internal/domain/blockchain"
	domain "go-blocker/internal/domain/payment"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/pkg/utils"
	"go-blocker/internal/rpc"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

type ETH struct {
	S *app.Service
}

func (w *ETH) Name() string {
	return "ETH"
}

func (w *ETH) Chain() string {
	return "ethereum"
}

func (w *ETH) Decimals() int {
	return 18
}

func (w *ETH) GetPendingBalance(client *ethclient.Client, wallet common.Address) big.Float {
	balance, err := client.PendingBalanceAt(context.Background(), wallet)
	if err != nil {
		logger.Log.Errorf("ETH: Error getting pending balance for address %s: %s", wallet.Hex(), err)
		return *big.NewFloat(0)
	}

	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))

	total := big.NewFloat(0).SetFloat64(0.06)
	if ethBalance.Cmp(total) > 0 {
		logger.Log.Infof("ETH: Pending balance for address %s: %f", wallet.Hex(), ethBalance)
		id, err := w.CheckAddress(wallet.Hex())
		if err != nil {
			return *ethBalance
		}
		w.S.Status(id, domain.Received, nil, nil, nil)
		storage.PaymentAddressStore.SetPending(wallet.Hex(), true)
	}
	return *ethBalance
}

func (w *ETH) HasActiveAddresses() bool {
	return storage.PaymentAddressStore.Len() > 0
}

func (w *ETH) CheckTransactions(m *rpc.Manager, client *ethclient.Client, block []*types.Receipt) (uuid.UUID, error) {
	if len(block) == 0 {
		return uuid.Nil, fmt.Errorf("no transactions in block")
	}

	blockbynum, err := client.BlockByNumber(context.Background(), block[0].BlockNumber)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error: %s", err)
	}

	isSent := false
	for _, tx := range blockbynum.Transactions() {
		if tx.To() == nil {
			continue
		}

		id, err := w.CheckAddress(tx.To().Hex())
		if err != nil {
			continue
		}

		var Tx *types.Transaction
		var eth, txid string
		Tx, eth, isSent = w.IsTransactionMatch(client, &app.CheckTxRequest{TxID: tx.Hash().Hex()})
		txid = Tx.Hash().Hex()

		logger.Log.Infof("ETH: Incoming transaction %s, Amount: %s", tx.Hash().Hex(), eth)
		w.S.Status(id, domain.Completed, &eth, &txid, nil)
	}

	if !isSent {
		w.CheckInternalTxs(m, client, block[0])
	}

	return uuid.Nil, fmt.Errorf("no matching address found")

}

func (w *ETH) CheckAddress(address string) (uuid.UUID, error) {
	// Base Transfer From - To
	if id, ok := storage.PaymentAddressStore.Get(address); ok {
		return id, nil
	}

	return uuid.Nil, fmt.Errorf("no matching address found")
}

func (w *ETH) CheckInternalTxs(m *rpc.Manager, client *ethclient.Client, block *types.Receipt) {
	_, url, err := m.GetClientForChain(blockchain.ChainType(w.Chain()))
	if err != nil {
		logger.Log.Debugf("[%s] No healthy RPC node for %s", w.Chain(), url)
		return
	}
	result, err := utils.TraceBlock(url, "0x"+block.BlockNumber.Text(16))
	if err != nil {
		logger.Log.Debugf("[%s] No healthy TraceBlock for %s: %s", w.Chain(), url, err)
		return
	}
	for _, tx := range result {
		id, err := w.CheckAddress(tx.Action.To)
		if err != nil {
			continue
		}

		wei := new(big.Int)
		wei.SetString(tx.Action.Value[2:], 16) // remove "0x" and parse as base 16

		eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18)).Text('f', 18)
		txid := tx.TransactionHash
		logger.Log.Infof("ETH: Incoming transaction %s, Amount: %s", txid, eth)

		IsStuck := true
		w.S.Status(id, domain.Completed, &eth, &txid, &IsStuck)
	}
}

func (w *ETH) IsTransactionMatch(client *ethclient.Client, tx *app.CheckTxRequest) (*types.Transaction, string, bool) {
	Tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.TxID))
	if err != nil {
		logger.Log.Debugf("ETH: %s", err)
	}

	wei := Tx.Value()
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18)).Text('f', 18)
	logger.Log.Infof("ETH: Incoming transaction %s, Amount: %s", Tx.Hash().Hex(), eth)

	return Tx, eth, true
}
