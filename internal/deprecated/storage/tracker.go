package storage

import "sync"

type BlockTracker interface {
	GetLastBlock(chain string) (uint64, error)
	SetLastBlock(chain string, block uint64) error
}

type MemoryTracker struct {
	mu     sync.RWMutex
	blocks map[string]uint64
}

func NewMemoryTracker() *MemoryTracker {
	return &MemoryTracker{
		blocks: make(map[string]uint64),
	}
}

func (t *MemoryTracker) GetLastBlock(chain string) (uint64, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.blocks[chain], nil
}

func (t *MemoryTracker) SetLastBlock(chain string, block uint64) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.blocks[chain] = block
	return nil
}
