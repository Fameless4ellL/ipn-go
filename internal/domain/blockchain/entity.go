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
