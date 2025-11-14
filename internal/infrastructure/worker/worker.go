package balancechecker

import (
	"context"
	"fmt"
	"log"
	"time"

	application "go-blocker/internal/application/payment"
	"go-blocker/internal/infrastructure/payment"
	logger "go-blocker/internal/pkg/log"
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

	for _, addr := range addresses {
		logger.Log.Debugf("Checking address: %s", addr.Address)

		if addr.Timeout.Before(time.Now()) {
			log.Printf("INFO: Skipping address %s due to timeout", addr.Address)
			utils.Send(map[string]interface{}{
				"status":          payment.Timeout,
				"address":         addr.Address,
				"stuck":           false,
				"received_amount": "0",
				"network":         string(addr.Network),
				"currency":        string(addr.Currency),
			}, addr.Callback)
			w.Service.Box.Delete(addr.Address)
			w.Service.Repo.Delete(addr.ID)
			continue
		}

		currency, err := w.Service.Provider.GetWatcher(addr.Network, addr.Currency)
		if err != nil {
			log.Printf("ERROR: No watcher for %s: %v", addr.Currency, err)
			w.Service.Box.Delete(addr.Address)
			w.Service.Repo.Delete(addr.ID)
			continue
		}

		isbalanced := currency.GetPendingBalance(addr.Address)
		if isbalanced {
			amount, isstuck := currency.GetLatestTx(addr.Address)
			if amount == "" {
				log.Printf("ERROR: No latest tx found for address %s", addr.Address)
				continue
			}
			utils.Send(map[string]interface{}{
				"status":          payment.Received,
				"address":         addr.Address,
				"stuck":           isstuck,
				"received_amount": fmt.Sprintf("%v", amount),
				"network":         string(addr.Network),
				"currency":        string(currency.GetName()),
			}, addr.Callback)
			w.Service.Box.Delete(addr.Address)
			w.Service.Repo.Delete(addr.ID)
		}

		time.Sleep(1 * time.Second) // Avoid rate limiting
	}
}
