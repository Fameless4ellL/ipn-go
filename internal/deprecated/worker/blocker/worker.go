package blocker

import (
	"go-blocker/internal/application/payment"
	"go-blocker/internal/deprecated/watcher"
	blockchain "go-blocker/internal/domain/blockchain"
	"go-blocker/internal/infrastructure/storage"
	"go-blocker/internal/rpc"
	"time"
)

func Start(
	service *payment.Service,
) {
	grouped := map[blockchain.ChainType][]watcher.CurrencyWatcher{
		blockchain.Ethereum: {
			// &watcher.ETH{S: service},
			// &watcher.USDT{S: service},
		},
	}

	manager := rpc.NewManager()
	manager.StartHealthMonitor(5 * time.Second)
	tracker := storage.NewMemoryTracker()

	go Pending(service, manager, grouped)
	go IncomingTx(service, manager, grouped, tracker)
}
