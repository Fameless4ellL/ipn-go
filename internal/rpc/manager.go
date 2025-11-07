package rpc

import (
	"log"
	"sync"
	"time"

	domain "go-blocker/internal/domain/blockchain"

	"github.com/ethereum/go-ethereum/ethclient"
)

type Manager struct {
	nodes []domain.RPCNode
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

func NewManager() *Manager {
	return &Manager{nodes: []domain.RPCNode{
		{URL: "https://eth.drpc.org", Chain: domain.Ethereum, Healthy: true},
		{URL: "https://api.mainnet-beta.solana.com", Chain: domain.Solana, Healthy: true},
		{URL: "https://bsc.drpc.org", Chain: domain.Binance, Healthy: true},
		{URL: "https://solana.drpc.org", Chain: domain.Solana, Healthy: true},
	}}
}

func (m *Manager) GetClientForChain(chain domain.ChainType) (*ethclient.Client, string, error) {
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
		log.Printf("RPC node failed will try to change another: %s (%v)", node.URL, err)
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
					log.Printf("RPC node recovered: %s", node.URL)
					m.nodes[i].Healthy = true
					m.nodes[i].LastFailure = time.Time{}
					client.Close()
				} else {
					log.Printf("RPC node still unhealthy: %s", node.URL)
				}
			}
			m.mu.Unlock()
		}
	}()
}
