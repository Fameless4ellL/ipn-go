package provider

import (
	"context"
	constants "go-blocker/internal/const"
	logger "go-blocker/internal/log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ETH struct{}

func (w *ETH) Name() string {
	return "ETH"
}

func (w *ETH) Chain() string {
	return "ethereum"
}

func (w *ETH) Decimals() int {
	return 18
}

func (w *ETH) IsTransactionMatch(client *ethclient.Client, tx *constants.CheckTxRequest) (string, bool) {
	Tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.TxID))
	if err != nil {
		logger.Log.Debugf("ETH: %s", err)
	}

	wei := Tx.Value()
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18)).Text('f', 18)
	logger.Log.Infof("ETH: Incoming transaction %s, Amount: %s", Tx.Hash().Hex(), eth)

	return eth, true
}
