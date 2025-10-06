package blockchain

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

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

type Address struct {
	ID       uuid.UUID
	Address  common.Address
	Currency CurrencyType
	Callback string
	Timeout  time.Time
}
