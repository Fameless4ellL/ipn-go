package blockchain

import "github.com/ethereum/go-ethereum/ethclient"

type Manager interface {
	GetClientForChain(chain ChainType) (*ethclient.Client, string, error)
}

type Watcher interface {
	GetWatcher(chain ChainType, currency CurrencyType) (Currency, error)
}

type Currency interface {
	Name() string
	Chain() string
	Decimals() int
	IsTransactionMatch(client *ethclient.Client, url string, address string, txid string) (string, bool)
	GetLatestTx(client *ethclient.Client, url string, address string) (string, bool)
}
