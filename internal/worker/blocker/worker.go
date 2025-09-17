package blocker

import (
	"go-blocker/internal/config"
	"go-blocker/internal/payment"
	"go-blocker/internal/rpc"
	"go-blocker/internal/storage"
	"go-blocker/internal/watcher"
	"time"
)

func Start(
	service *payment.PaymentService,
) {
	grouped := map[rpc.ChainType][]watcher.CurrencyWatcher{
		rpc.Ethereum: {
			&watcher.ETH{S: service},
			&watcher.USDT{S: service},
		},
	}

	manager := rpc.NewManager(config.Nodes)
	manager.StartHealthMonitor(5 * time.Second)
	tracker := storage.NewMemoryTracker()

	go Pending(service, manager, grouped)
	go IncomingTx(service, manager, grouped, tracker)
}
