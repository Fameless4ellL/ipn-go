package blockchain

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

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
	GetPendingBalance(client *ethclient.Client, wallet common.Address) bool
	IsTransactionMatch(client *ethclient.Client, url string, address string, txid string) (string, bool)
	GetLatestTx(client *ethclient.Client, url string, address string) (string, bool)
}

type Storage interface {
	List() []Address
	Len() int
	Set(address string, id uuid.UUID, c CurrencyType, callback string, timeout time.Time)
	Get(string) (Address, bool)
	Delete(address string)
}
