package watcher

import (
	"go-blocker/internal/rpc"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

type CurrencyWatcher interface {
	Name() string
	Chain() string
	Decimals() int
	HasActiveAddresses() bool
	CheckTransactions(m *rpc.Manager, client *ethclient.Client, block []*types.Receipt) (uuid.UUID, error)
	GetPendingBalance(client *ethclient.Client, wallet common.Address) big.Float
}
