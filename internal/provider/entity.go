package provider

import (
	constants "go-blocker/internal/const"

	"github.com/ethereum/go-ethereum/ethclient"
)

type CurrencyWatcher interface {
	Name() constants.CurrencyType
	Chain() string
	Decimals() int
	IsTransactionMatch(client *ethclient.Client, url string, tx *constants.CheckTxRequest) (string, bool)
	GetLatestTx(client *ethclient.Client, url string, req constants.FindTxRequest) (string, bool)
}
