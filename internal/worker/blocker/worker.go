package blocker

import (
	"go-blocker/internal/payment"
	"go-blocker/internal/rpc"
	"go-blocker/internal/storage"
	"go-blocker/internal/watcher"
	"time"
)

func Start(
	service *payment.PaymentService,
) {
	nodes := []rpc.RPCNode{
		// {URL: "https://eth.drpc.org", Chain: rpc.Ethereum, Healthy: true, Processing: false}, // has trace_block
		// {URL: "https://api.noderpc.xyz/rpc-mainnet/public", Chain: rpc.Ethereum, Healthy: true, Processing: false}, // has trace_block
		// {URL: "https://ethereum-public.nodies.app", Chain: rpc.Ethereum, Healthy: true, Processing: false}, // has trace_block
		// {URL: "https://endpoints.omniatech.io/v1/eth/mainnet/public", Chain: rpc.Ethereum, Healthy: true, Processing: false}, // has trace_block
		// {URL: "https://eth.api.onfinality.io/public", Chain: rpc.Ethereum, Healthy: true, Processing: false}, // has trace_block
		// {URL: "https://eth.llamarpc.com", Chain: rpc.Ethereum, Healthy: true}, not trace_block
		// {URL: "https://ethereum-rpc.publicnode.com", Chain: rpc.Ethereum, Healthy: true}, // not trace_block
		// {URL: "https://go.getblock.io/aefd01aa907c4805ba3c00a9e5b48c6b", Chain: rpc.Ethereum, Healthy: true}, too many requests and no support for trace_block
		{URL: "https://sepolia.drpc.org", Chain: rpc.Ethereum, Healthy: true}, // test
	}
	grouped := map[rpc.ChainType][]watcher.CurrencyWatcher{
		rpc.Ethereum: {
			&watcher.ETH{S: service},
			&watcher.USDT{S: service},
		},
	}

	manager := rpc.NewManager(nodes)
	manager.StartHealthMonitor(5 * time.Second)
	tracker := storage.NewMemoryTracker()

	go Pending(service, manager, grouped)
	go IncomingTx(service, manager, grouped, tracker)
}
