package storage

import (
	"strings"
	"sync"

	"github.com/google/uuid"
)

type Address struct {
	ID      uuid.UUID
	Pending bool
}

type AddressStore struct {
	mu    sync.RWMutex
	store map[string]Address
}

var PaymentAddressStore *AddressStore

func InitStores() {
	PaymentAddressStore = NewAddressStore()
}

func NewAddressStore() *AddressStore {
	return &AddressStore{
		store: make(map[string]Address),
	}
}

func (s *AddressStore) Set(address string, id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[strings.ToLower(address)] = Address{ID: id, Pending: false}
}

func (s *AddressStore) Get(address string) (uuid.UUID, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.store[strings.ToLower(address)]
	return id.ID, ok
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

func (s *AddressStore) List() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	addresses := make([]string, 0, len(s.store))
	for addr := range s.store {
		addresses = append(addresses, strings.ToLower(addr))
	}
	return addresses
}

func (s *AddressStore) SetPending(address string, isPending bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	addr, exists := s.store[strings.ToLower(address)]
	if exists {
		addr.Pending = isPending
		s.store[strings.ToLower(address)] = addr
	}
}

func (s *AddressStore) IsPending(address string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	addr, exists := s.store[strings.ToLower(address)]
	if exists {
		return addr.Pending
	}
	return false
}
