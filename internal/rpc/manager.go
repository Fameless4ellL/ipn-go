package rpc

import (
	"go-blocker/internal/config"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Manager struct {
	nodes []RPCNode
	index int
	mu    sync.Mutex
}

var ErrAllNodesFailed = &RPCError{"All RPC nodes failed"}

type RPCError struct {
	Msg string
}

func (e *RPCError) Error() string {
	return e.Msg
}

func NewManager(nodes []RPCNode) *Manager {
	return &Manager{nodes: nodes}
}

func (m *Manager) GetClientForChain(chain ChainType) (*ethclient.Client, string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i := 0; i < len(m.nodes); i++ {
		node := m.nodes[m.index]
		m.index = (m.index + 1) % len(m.nodes)

		if node.Chain != chain || !node.Healthy {
			continue
		}

		client, err := ethclient.Dial(node.URL)
		if err == nil {
			return client, node.URL, nil
		}
		config.Log.Infof("RPC node failed will try to change another: %s (%v)", node.URL, err)
		m.nodes[m.index].Healthy = false
		m.nodes[m.index].LastFailure = time.Now()
	}

	return nil, "", ErrAllNodesFailed
}

func (m *Manager) StartHealthMonitor(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			m.mu.Lock()
			for i, node := range m.nodes {
				if node.Healthy {
					continue
				}

				client, err := ethclient.Dial(node.URL)
				if err == nil {
					config.Log.Infof("RPC node recovered: %s", node.URL)
					m.nodes[i].Healthy = true
					m.nodes[i].LastFailure = time.Time{}
					client.Close()
				} else {
					config.Log.Infof("RPC node still unhealthy: %s", node.URL)
				}
			}
			m.mu.Unlock()
		}
	}()
}
