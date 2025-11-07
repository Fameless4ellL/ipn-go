package provider

import (
	"fmt"
	"go-blocker/internal/domain/blockchain"
	"go-blocker/internal/infrastructure/provider/etherscan"

	"github.com/ethereum/go-ethereum/common"
)

type Registry struct {
	// Keep the map private (lowercase 'w') so access is forced through the method.
	watchers map[blockchain.ChainType]blockchain.Chain
}

func NewCurrencyWatcherRegistry(m blockchain.Manager) *Registry {
	client, url, err := m.GetClientForChain(blockchain.Ethereum)
	if err != nil {
		panic(fmt.Sprintf("Failed to create RPC client: %v", err))
	}

	evm := &EVM{
		client: client,
		scan:   etherscan.NewClient("1"),
		url:    url,
	}

	client, url, err = m.GetClientForChain(blockchain.Binance)
	if err != nil {
		panic(fmt.Sprintf("Failed to create RPC client: %v", err))
	}

	bsc := &EVM{
		client: client,
		scan:   etherscan.NewClient("56"),
		url:    url,
	}

	return &Registry{
		watchers: map[blockchain.ChainType]blockchain.Chain{
			blockchain.Ethereum: {
				Id:   1,
				Name: blockchain.Ethereum,
				Currencies: map[blockchain.CurrencyType]blockchain.Currency{
					blockchain.ETH: &Native[blockchain.API]{
						client:   evm,
						Name:     blockchain.ETH,
						Decimals: 18,
					},
					blockchain.USDT: &ERC20[blockchain.API]{
						client:   evm,
						Name:     blockchain.USDT,
						Address:  common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
						Decimals: 6,
					},
					blockchain.USDC: &ERC20[blockchain.API]{
						client:   evm,
						Name:     blockchain.USDC,
						Address:  common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"),
						Decimals: 6,
					},
				},
			},
			blockchain.Binance: {
				Id:   56,
				Name: blockchain.Binance,
				Currencies: map[blockchain.CurrencyType]blockchain.Currency{
					blockchain.BNB: &Native[blockchain.API]{
						client:   bsc,
						Name:     blockchain.BNB,
						Decimals: 18,
					},
					blockchain.USDC: &ERC20[blockchain.API]{
						client:   bsc,
						Name:     blockchain.USDC,
						Address:  common.HexToAddress("0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d"),
						Decimals: 18,
					},
					blockchain.BUSD: &ERC20[blockchain.API]{
						client:   bsc,
						Name:     blockchain.BUSD,
						Address:  common.HexToAddress("0x55d398326f99059fF775485246999027B3197955"),
						Decimals: 18,
					},
				},
			},
		},
	}
}

// GetWatcher retrieves the CurrencyWatcher for the given chain and currency.
func (r *Registry) GetWatcher(chain blockchain.ChainType, currency blockchain.CurrencyType) (blockchain.Currency, error) {
	c, ok := r.watchers[chain]
	if !ok {
		return nil, fmt.Errorf("[%s] is not supported by the registry", chain)
	}

	watcher, ok := c.Currencies[currency]
	if !ok {
		return nil, fmt.Errorf("[%s] not supported on chain '%s'", currency, chain)
	}

	return watcher, nil
}
