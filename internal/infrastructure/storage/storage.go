package storage

import (
	"strings"
	"sync"
	"time"

	blockchain "go-blocker/internal/domain/blockchain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type AddressStore struct {
	mu    sync.RWMutex
	store map[string]blockchain.Address
}

var PaymentAddressStore *AddressStore

func InitStores() {
	PaymentAddressStore = NewAddressStore()
}

func NewAddressStore() *AddressStore {
	return &AddressStore{
		store: make(map[string]blockchain.Address),
	}
}

func (s *AddressStore) Set(
	address string,
	id uuid.UUID,
	c blockchain.CurrencyType,
	callback string,
	timeout time.Time,
) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[strings.ToLower(address)] = blockchain.Address{
		ID:       id,
		Currency: c,
		Address:  common.HexToAddress(address),
		Callback: callback,
		Timeout:  timeout,
	}
}

func (s *AddressStore) Get(address string) (blockchain.Address, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.store[strings.ToLower(address)]
	return id, ok
}

func (s *AddressStore) Delete(address string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, strings.ToLower(address))
}

func (s *AddressStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.store)
}

func (s *AddressStore) List() []blockchain.Address {
	s.mu.RLock()
	defer s.mu.RUnlock()
	addresses := make([]blockchain.Address, 0, len(s.store))
	for _, addr := range s.store {
		addresses = append(addresses, addr)
	}
	return addresses
}
