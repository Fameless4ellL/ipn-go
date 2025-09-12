package blocker

import (
	"context"
	"fmt"
	"go-blocker/internal/config"
	"go-blocker/internal/payment"
	"go-blocker/internal/rpc"
	"go-blocker/internal/storage"
	"go-blocker/internal/watcher"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	w3 "github.com/ethereum/go-ethereum/rpc"

	"time"
)

var (
	processingMutex sync.Mutex
	processing      = make(map[rpc.ChainType]bool)
)

func IncomingTx(
	s *payment.PaymentService,
	manager *rpc.Manager,
	watchersmaped map[rpc.ChainType][]watcher.CurrencyWatcher,
	tracker storage.BlockTracker,
) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for chain, watchers := range watchersmaped {
			if storage.PaymentAddressStore.Len() == 0 {
				continue
			}

			processingMutex.Lock()
			if processing[chain] {
				processingMutex.Unlock()
				continue
			}
			processing[chain] = true
			processingMutex.Unlock()

			client, url, err := manager.GetClientForChain(rpc.ChainType(chain))
			if err != nil {
				config.Log.Debugf("[%s] No healthy RPC nodes for %s", chain, url)
				processingMutex.Lock()
				processing[chain] = false
				processingMutex.Unlock()
				continue
			}
			latest, err := client.BlockNumber(context.Background())
			if err != nil {
				config.Log.Debugf("[%s] Failed to get latest block from %s: %v", chain, url, err)
				processingMutex.Lock()
				processing[chain] = false
				processingMutex.Unlock()
				continue
			}
			go func(chain rpc.ChainType, watchers []watcher.CurrencyWatcher, client *ethclient.Client, latest uint64) {
				defer func() {
					processingMutex.Lock()
					processing[chain] = false
					processingMutex.Unlock()
				}()

				err := Blocker(s, manager, watchers, client, latest, tracker)
				if err != nil {
					config.Log.Errorf("[%s] Blocker error: %v", chain, err)
				}

				time.Sleep(1 * time.Second)
			}(chain, watchers, client, latest)

		}
	}
}

func Blocker(
	s *payment.PaymentService,
	m *rpc.Manager,
	watchers []watcher.CurrencyWatcher,
	client *ethclient.Client,
	latest uint64,
	tracker storage.BlockTracker,
) error {
	chain := watchers[0].Chain()

	lastBlock, err := tracker.GetLastBlock(string(chain))
	if err != nil || lastBlock == 0 || latest-lastBlock > 10 {
		lastBlock = latest - 10 // fallback
	}

	config.Log.Infof("[%s] Looking between  %d | %d", chain, lastBlock, latest)

	for blockNum := lastBlock + 1; blockNum <= latest; blockNum++ {
		blockRef := w3.BlockNumberOrHashWithNumber(w3.BlockNumber(blockNum))
		receipts, err := client.BlockReceipts(context.Background(), blockRef)
		if err != nil {
			config.Log.Debugf("[%s] Failed to get receipts for block %d: %v", chain, blockNum, err)
			return fmt.Errorf("[%s] Failed to get receipts for block %d: %v", chain, blockNum, err)
		}

		for _, w := range watchers {
			w.CheckTransactions(m, client, receipts)
		}
		tracker.SetLastBlock(string(chain), blockNum)
	}
	return nil
}
