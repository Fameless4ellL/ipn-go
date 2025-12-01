package solana

import (
	"fmt"
	"testing"
)

func TestRPCGetBalance(t *testing.T) {
	client := NewRPC("https://api.mainnet-beta.solana.com")
	address := "83astBRguLMdt2h5U1Tpdq5tjFoJ6noeGwaY3mDLVcri"

	balance, err := client.GetBalance(address)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if balance == nil {
		t.Fatalf("expected balance, got empty string")
	}
}

func TestRPCGetTokenAccountBalance(t *testing.T) {
	client := NewRPC("https://api.mainnet-beta.solana.com")
	address := "417wUmxD6mPUKfrdeFgUveLfWa9q1LkJt1S5Dz9BxHqS"

	balance, err := client.GetTokenAccountBalance(address)
	fmt.Printf("balance: %v\n", balance)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if balance == nil {
		t.Fatalf("expected balance, got empty string")
	}
}
