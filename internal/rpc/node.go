package rpc

import "time"

type ChainType string

const (
	Ethereum ChainType = "ethereum"
)

type RPCNode struct {
	URL         string
	Chain       ChainType
	LastFailure time.Time
	Healthy     bool
}
