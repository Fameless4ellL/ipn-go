package provider

import (
	"fmt"
	blockchain "go-blocker/internal/domain/blockchain"
)

type CurrencyWatcherRegistry struct {
	// Keep the map private (lowercase 'w') so access is forced through the method.
	watchers map[blockchain.ChainType]map[blockchain.CurrencyType]blockchain.Currency
}

// NewCurrencyWatcherRegistry initializes and returns the registry instance.
// This is the single place where concrete implementations are wired.
func NewCurrencyWatcherRegistry() *CurrencyWatcherRegistry {
	return &CurrencyWatcherRegistry{
		watchers: map[blockchain.ChainType]map[blockchain.CurrencyType]blockchain.Currency{
			blockchain.Ethereum: {
				blockchain.ETH:  &ETH{},
				blockchain.USDT: &USDT{},
				blockchain.USDC: &USDC{},
			},
			// Add other chains here...
		},
	}
}

// GetWatcher retrieves the CurrencyWatcher for the given chain and currency.
func (r *CurrencyWatcherRegistry) GetWatcher(chain blockchain.ChainType, currency blockchain.CurrencyType) (blockchain.Currency, error) {
	currencyMap, chainFound := r.watchers[chain]
	if !chainFound {
		return nil, fmt.Errorf("chain type '%s' is not supported by the registry", chain)
	}

	watcher, currencyFound := currencyMap[currency]
	if !currencyFound {
		return nil, fmt.Errorf("currency type '%s' not supported on chain '%s'", currency, chain)
	}

	return watcher, nil
}
