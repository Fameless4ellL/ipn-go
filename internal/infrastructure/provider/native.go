package provider

import (
	"go-blocker/internal/domain/blockchain"
	logger "go-blocker/internal/pkg/log"
	"log/slog"
	"math/big"
	"strings"
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

func (w *Native[T]) CheckInternalTxs(Tx *blockchain.Transaction, address string) (string, bool) {
	blocknumberInt := new(big.Int)
	blocknumberInt.SetString(Tx.BlockNumber.String(), 10)
	blocknumberHex := blocknumberInt.Text(16)

	result, err := w.client.TraceBlock(blocknumberHex, address)
	if err != nil {
		logger.Log.Error(
			"Error",
			slog.Group(string(w.Name),
				slog.String("blocknumber", blocknumberHex),
				slog.String("address", address),
				slog.Any("error", err),
			),
		)
		return "", false
	}
	for _, tx := range result {
		if !strings.EqualFold(Tx.Hash, tx.TransactionHash) {
			slog.Debug("tx hash mismatch",
				slog.String("tx.Hash", Tx.Hash),
				slog.String("tx.TransactionHash", tx.TransactionHash),
			)
			continue
		}
		if !strings.EqualFold(tx.Action.To, address) {
			slog.Debug("tx to address mismatch",
				slog.String("tx.Action.To", tx.Action.To),
				slog.String("address", address),
			)
			continue
		}
		wei := new(big.Int)
		wei.SetString(tx.Action.Value[2:], 16) // remove "0x" and parse as base 16
		weiFloat := new(big.Float).SetInt(wei)
		ethValue := new(big.Float).Quo(weiFloat, big.NewFloat(1e18))
		ethStr := ethValue.Text('f', 18)
		txid := tx.TransactionHash
		logger.Log.Info(
			"Incoming Internal transaction",
			slog.Group(string(w.Name),
				slog.String("txid", txid),
				slog.String("amount", ethStr),
			),
		)

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

	if !strings.EqualFold(Tx.ContractAddress, address) {
		Tx, err := w.client.TransactionReceipt(txid)
		if err != nil {
			logger.Log.Warn(
				"Error getting transaction receipt",
				slog.Group(string(w.Name),
					slog.String("txid", txid),
					slog.Any("error", err),
				),
			)
			return "", false
		}
		return w.CheckInternalTxs(Tx, address)
	}

	eth := new(big.Float).Quo(new(big.Float).SetInt(Tx.Value), big.NewFloat(1e18)).Text('f', 18)
	logger.Log.Info(
		"Incoming transaction",
		slog.Group(string(w.Name),
			slog.String("txid", Tx.Hash),
			slog.String("amount", eth),
		),
	)

	return eth, false
}

func (w *Native[T]) GetLatestTx(address string) (string, bool) {
	hash, err := w.client.GetTx(address)
	if err != nil {
		logger.Log.Warn(
			"Internal Error",
			slog.Group(string(w.Name),
				slog.Any("error", err),
			),
		)
	}

	return w.IsTransactionMatch(address, hash)
}
