package blockchain

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"resty.dev/v3"
)

type ChainType string
type CurrencyType string

const (
	Ethereum ChainType = "ethereum"
	Binance  ChainType = "binance"
	Solana   ChainType = "solana"
	Litecoin ChainType = "litecoin"
)

const (
	ETH  CurrencyType = "ETH"
	BNB  CurrencyType = "BNB"
	USDT CurrencyType = "USDT"
	BUSD CurrencyType = "BUSD"
	USDC CurrencyType = "USDC"
	SOL  CurrencyType = "SOL"
	LTC  CurrencyType = "LTC"
)

type Chain struct {
	Id         int
	Name       ChainType
	Currencies map[CurrencyType]Currency
}

type Address struct {
	ID       uuid.UUID
	Address  string
	Network  ChainType
	Currency CurrencyType
	Callback string
	Timeout  time.Time
	Amount   string // expected amount
}

// MinStuckPercent is the minimal share (in percent) of the expected amount
// that a stuck transaction must carry to be accepted.
const MinStuckPercent = 5

// IsStuckAmountEnough reports whether received is at least MinStuckPercent of
// the expected amount. If the expected amount is unknown, the check is skipped.
func (a Address) IsStuckAmountEnough(received string) bool {
	expected, ok := new(big.Rat).SetString(a.Amount)
	if !ok || expected.Sign() <= 0 {
		return true
	}
	got, ok := new(big.Rat).SetString(received)
	if !ok {
		return false
	}
	min := new(big.Rat).Mul(expected, big.NewRat(MinStuckPercent, 100))
	return got.Cmp(min) >= 0
}

type Logs struct {
	Address string
	Topics  []common.Hash
	Data    []byte
}

type Transaction struct {
	BlockNumber     *big.Int
	ContractAddress string
	Hash            string
	Logs            []*Logs
	Value           *big.Int
}

type Node struct {
	Rr     *resty.RoundRobin
	Client *ethclient.Client
}
