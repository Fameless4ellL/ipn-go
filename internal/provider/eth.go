package provider

import (
	"context"
	constants "go-blocker/internal/const"
	logger "go-blocker/internal/log"
	"go-blocker/internal/provider/etherscan"
	"go-blocker/internal/utils"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ETH struct{}

func (w *ETH) Name() constants.CurrencyType {
	return constants.ETH
}

func (w *ETH) Chain() string {
	return "ethereum"
}

func (w *ETH) Decimals() int {
	return 18
}

func (w *ETH) CheckInternalTxs(url string, Tx *types.Receipt, address string) (string, bool) {
	blocknumberInt := new(big.Int)
	blocknumberInt.SetString(Tx.BlockNumber.String(), 10)
	blocknumberHex := blocknumberInt.Text(16)
	result, err := utils.TraceBlock(url, "0x"+blocknumberHex)
	if err != nil {
		logger.Log.Debugf("[%s] No healthy TraceBlock: %s", w.Chain(), err)
		return "", false
	}
	for _, tx := range result {
		if Tx.TxHash.Hex() != tx.TransactionHash {
			continue
		}
		if !strings.EqualFold(tx.Action.To, address) {
			continue
		}
		wei := new(big.Int)
		wei.SetString(tx.Action.Value[2:], 16) // remove "0x" and parse as base 16
		weiFloat := new(big.Float).SetInt(wei)
		ethValue := new(big.Float).Quo(weiFloat, big.NewFloat(1e18))
		ethStr := ethValue.Text('f', 18)
		txid := tx.TransactionHash
		logger.Log.Infof("ETH: Incoming Internal transaction %s, Amount: %s", txid, ethStr)

		IsStuck := true
		return ethStr, IsStuck
	}
	return "", false
}

func (w *ETH) IsTransactionMatch(client *ethclient.Client, url string, tx *constants.CheckTxRequest) (string, bool) {
	Tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.TxID))
	if err != nil {
		return "", false
	}

	if *Tx.To() != common.HexToAddress(tx.Address) {
		Tx, err := client.TransactionReceipt(context.Background(), common.HexToHash(tx.TxID))
		if err != nil {
			return "", false
		}
		return w.CheckInternalTxs(url, Tx, tx.Address)
	}

	wei := Tx.Value()
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18)).Text('f', 18)

	logger.Log.Infof("ETH: Incoming transaction %s, Amount: %s", Tx.Hash().Hex(), eth)

	return eth, true
}

func (w *ETH) GetLatestTx(client *ethclient.Client, url string, req constants.FindTxRequest) (string, bool) {
	ScanClient := etherscan.NewClient()
	resp, err := ScanClient.GetTransactions(req.Address, 0, 99999999, 1, 10, "asc")
	if err != nil {
		logger.Log.Warnf("Error: %v", err)
		return "", false
	}

	return w.IsTransactionMatch(client, url, &constants.CheckTxRequest{
		Address:  req.Address,
		Currency: w.Name(),
		TxID:     resp.Result[0].Hash,
	})
}
