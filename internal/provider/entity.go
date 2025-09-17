package provider

import (
	constants "go-blocker/internal/const"

	"github.com/ethereum/go-ethereum/ethclient"
)

type CurrencyWatcher interface {
	Name() string
	Chain() string
	Decimals() int
	IsTransactionMatch(client *ethclient.Client, tx *constants.CheckTxRequest) (string, bool)
}
