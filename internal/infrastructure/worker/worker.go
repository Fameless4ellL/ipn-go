package balancechecker

import (
	"context"
	"fmt"
	"log"
	"time"

	application "go-blocker/internal/application/payment"
	blockchain "go-blocker/internal/domain/blockchain"
	"go-blocker/internal/infrastructure/payment"
	"go-blocker/internal/pkg/utils"
)

type Worker struct {
	Service  *application.Service
	Interval time.Duration
}

func NewWorker(s *application.Service, interval time.Duration) *Worker {
	return &Worker{
		Service:  s,
		Interval: interval,
	}
}

// Start runs the periodic task loop.
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()

	log.Println("Address tracker started.")
	w.executeCheck() // Initial run

	for {
		select {
		case <-ctx.Done():
			log.Println("Address tracker stopped.")
			return
		case <-ticker.C:
			w.executeCheck()
		}
	}
}

func (w *Worker) executeCheck() {
	addresses := w.Service.Box.List()

	client, url, err := w.Service.Manager.GetClientForChain(blockchain.Ethereum)
	if err != nil {
		log.Printf("ERROR: No healthy RPC nodes for Ethereum: %v", err)
		return
	}

	for _, addr := range addresses {
		if addr.Timeout.Before(time.Now()) {
			log.Printf("INFO: Skipping address %s due to timeout", addr.Address.String())
			utils.Send(map[string]interface{}{
				"status":          payment.Timeout,
				"address":         addr.Address.String(),
				"stuck":           true,
				"received_amount": "0",
				// "txid":            "",
				"currency": string(addr.Currency),
			}, addr.Callback)
			w.Service.Box.Delete(addr.Address.String())
			continue
		}

		currency, err := w.Service.Provider.GetWatcher(blockchain.Ethereum, addr.Currency)
		if err != nil {
			log.Printf("ERROR: No watcher for currency %s: %v", addr.Currency, err)
			w.Service.Box.Delete(addr.Address.String())
			continue
		}
		isbalanced := currency.GetPendingBalance(client, addr.Address)
		if isbalanced {
			amount, isstuck := currency.GetLatestTx(client, url, addr.Address.String())
			if amount == "" {
				log.Printf("ERROR: No latest tx found for address %s", addr.Address.String())
				continue
			}
			utils.Send(map[string]interface{}{
				"status":          payment.Received,
				"address":         addr.Address.String(),
				"stuck":           isstuck,
				"received_amount": fmt.Sprintf("%v", amount),
				"currency":        string(currency.Name()),
			}, addr.Callback)
			w.Service.Box.Delete(addr.Address.String())
		}

		time.Sleep(1 * time.Second) // Avoid rate limiting
	}
}
