package blockchain

import "time"

type ChainType string
type CurrencyType string

const (
	Ethereum ChainType = "ethereum"
)

const (
	ETH  CurrencyType = "ETH"
	USDT CurrencyType = "USDT"
	USDC CurrencyType = "USDC"
)

type RPCNode struct {
	URL         string
	Chain       ChainType
	LastFailure time.Time
	Healthy     bool
}