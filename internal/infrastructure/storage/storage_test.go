package storage

import (
    "sync"
    "testing"
    "time"

    "github.com/google/uuid"
    blockchain "go-blocker/internal/domain/blockchain"
)

func TestAddressStore_SetGetDeleteList(t *testing.T) {
    s := NewAddressStore()

    id := uuid.New()
    addr := "0xAbCdEf0123456789"
    now := time.Now().Add(1 * time.Hour)

    s.Set(addr, id, blockchain.Ethereum, blockchain.ETH, "http://cb", now)

    got, ok := s.Get(addr)
    if !ok {
        t.Fatalf("expected address to exist")
    }
    if got.ID != id {
        t.Fatalf("expected ID %v, got %v", id, got.ID)
    }
    if got.Callback != "http://cb" {
        t.Fatalf("expected callback %s, got %s", "http://cb", got.Callback)
    }

    list := s.List()
    if len(list) != 1 {
        t.Fatalf("expected list length 1, got %d", len(list))
    }

    s.Delete(addr)
    _, ok = s.Get(addr)
    if ok {
        t.Fatalf("expected address to be deleted")
    }
}

func TestAddressStore_ConcurrentAccess(t *testing.T) {
    s := NewAddressStore()
    var wg sync.WaitGroup
    const n = 200

    // writers
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            id := uuid.New()
            a := uuid.New().String()
            s.Set(a, id, blockchain.Ethereum, blockchain.ETH, "", time.Now().Add(time.Minute))
        }(i)
    }

    // readers
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            _ = s.List()
        }()
    }

    wg.Wait()
    // basic sanity: no panic/race; ensure map length > 0
    if len(s.List()) == 0 {
        t.Fatalf("expected non-empty store after concurrent writes")
    }
}