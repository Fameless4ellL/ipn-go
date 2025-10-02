package provider

import "github.com/ethereum/go-ethereum/ethclient"

type CurrencyWatcher interface {
	Name() string
	Chain() string
	Decimals() int
	IsTransactionMatch(client *ethclient.Client, url string, address string, txid string) (string, bool)
	GetLatestTx(client *ethclient.Client, url string, address string) (string, bool)
}
