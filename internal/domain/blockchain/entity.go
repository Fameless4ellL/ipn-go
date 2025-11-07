package blockchain

import (
	"time"

	"github.com/google/uuid"
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

type RPCNode struct {
	URL         string
	Chain       ChainType
	LastFailure time.Time
	Healthy     bool
}

type Address struct {
	ID       uuid.UUID
	Address  string
	Network  ChainType
	Currency CurrencyType
	Callback string
	Timeout  time.Time
}
