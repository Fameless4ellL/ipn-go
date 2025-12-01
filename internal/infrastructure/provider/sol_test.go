package provider

import (
	"go-blocker/internal/infrastructure/provider/solana"
	"testing"
)

func TestSolTransactionByHash(t *testing.T) {
	// Test cases for SolTransactionByHash
	sol := &SOL{
		client: solana.NewRPC("https://api.mainnet-beta.solana.com"),
		scan:   solana.NewClient(),
	}
	tx, err := sol.TransactionByHash("4mdSjU3NG5vophMgEpYT7WNzBszsARbeEg34vijKsYyj7ERUjhkxnYAEQLbuJmvZMjM67KyPhzLZw767fuJ5Z9EV")
	if err != nil {
		t.Errorf("Error fetching transaction: %v", err)
	}
	if tx == nil {
		t.Errorf("Transaction is nil")
	}
}
