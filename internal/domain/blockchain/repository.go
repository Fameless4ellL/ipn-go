package blockchain

import (
	"go-blocker/internal/pkg/utils"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
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
	GetName() CurrencyType
	GetPendingBalance(wallet string) bool
	IsTransactionMatch(address, txid string) (string, bool)
	GetLatestTx(address string) (string, bool)
}

type API interface {
	GetERC20Balance(abi, contract, wallet string) *big.Int
	GetERC20(contract, address string) (string, error)
	GetTx(address string) (string, error)
	GetBalance(wallet string) *big.Int
	TraceBlock(Tx *types.Receipt, address string) ([]utils.TraceResult, error)
	TransactionByHash(txid string) (*types.Transaction, error)
	TransactionReceipt(txid string) (*types.Receipt, error)
}

type Storage interface {
	List() []Address
	Len() int
	Set(address string, id uuid.UUID, n ChainType, c CurrencyType, callback string, timeout time.Time)
	Get(string) (Address, bool)
	Delete(address string)
}
