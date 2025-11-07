package provider

import (
	"go-blocker/internal/domain/blockchain"
	logger "go-blocker/internal/pkg/log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Native[T blockchain.API] struct {
	client   T
	Name     blockchain.CurrencyType
	Decimals int
}

func (w *Native[T]) GetName() blockchain.CurrencyType {
	return w.Name
}

func (w *Native[T]) GetPendingBalance(wallet string) bool {
	balance := w.client.GetBalance(wallet)
	Balance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	total := big.NewFloat(0)
	return Balance.Cmp(total) > 0
}

func (w *Native[T]) CheckInternalTxs(Tx *types.Receipt, address string) (string, bool) {
	result, err := w.client.TraceBlock(Tx, address)
	if err != nil {
		logger.Log.Debugf("[%s]: %s", w.Name, err)
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
		logger.Log.Infof("[%s]: Incoming Internal transaction %s, Amount: %s", w.Name, txid, ethStr)

		IsStuck := true
		return ethStr, IsStuck
	}
	return "", false
}

func (w *Native[T]) IsTransactionMatch(address string, txid string) (string, bool) {
	Tx, err := w.client.TransactionByHash(txid)
	if err != nil {
		return "", false
	}

	if *Tx.To() != common.HexToAddress(address) {
		Tx, err := w.client.TransactionReceipt(txid)
		if err != nil {
			logger.Log.Warnf("[%s]: Error getting transaction receipt for tx %s: %s", txid, err)
			return "", false
		}
		return w.CheckInternalTxs(Tx, address)
	}

	wei := Tx.Value()
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18)).Text('f', 18)

	logger.Log.Infof("[%s]: Incoming transaction %s, Amount: %s", w.Name, Tx.Hash().Hex(), eth)

	return eth, false
}

func (w *Native[T]) GetLatestTx(address string) (string, bool) {
	hash, err := w.client.GetTx(address)
	if err != nil {
		logger.Log.Warnf("Internal Error: %v", err)
	}

	return w.IsTransactionMatch(address, hash)
}
