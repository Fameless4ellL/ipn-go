package blocker

import (
	"go-blocker/internal/application/payment"
	"go-blocker/internal/deprecated/storage"
	"go-blocker/internal/deprecated/watcher"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/rpc"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	pendingMutex sync.Mutex
	pending      = make(map[rpc.ChainType]bool)
)

func Pending(
	s *payment.Service,
	manager *rpc.Manager,
	watchersmaped map[rpc.ChainType][]watcher.CurrencyWatcher,
) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		for chain, watchers := range watchersmaped {
			if storage.PaymentAddressStore.Len() == 0 {
				continue
			}

			pendingMutex.Lock()
			if pending[chain] {
				processingMutex.Unlock()
				continue
			}
			pending[chain] = true
			pendingMutex.Unlock()

			client, url, err := manager.GetClientForChain(rpc.ChainType(chain))
			if err != nil {
				logger.Log.Debugf("[%s] No healthy RPC nodes for %s", chain, url)
				pendingMutex.Lock()
				pending[chain] = false
				pendingMutex.Unlock()
				continue
			}
			go func(chain rpc.ChainType, watchers []watcher.CurrencyWatcher, client *ethclient.Client) {
				defer func() {
					pendingMutex.Lock()
					pending[chain] = false
					pendingMutex.Unlock()
				}()
				for _, address := range storage.PaymentAddressStore.List() {
					if storage.PaymentAddressStore.IsPending(address) {
						continue
					}

					for _, w := range watchers {
						w.GetPendingBalance(client, common.HexToAddress(address))
					}
				}

				time.Sleep(1 * time.Second)
			}(chain, watchers, client)

		}
	}
}
