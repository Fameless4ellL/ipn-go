package worker

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	application "go-blocker/internal/application/payment"
	"go-blocker/internal/domain/blockchain"
	"go-blocker/internal/infrastructure/payment"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/pkg/utils"
)

type Result struct {
	Data   *map[string]any
	Delete bool
}

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

	logger.Log.Info("Address tracker started.")

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Address tracker stopped.")
			return
		case <-ticker.C:
			w.checkAllAdress()
		}
	}
}


// Start to check all addresses in the Service box.
func (w *Worker) checkAllAdress() {
	addresses := w.Service.Box.List()
	for _, addr := range addresses {
		result := w.check(addr)
		if result == nil {
			continue
		}

		if result.Data != nil {
			utils.Send(map[string]any{
				"address":  addr.Address,
				"network":  string(addr.Network),
				"currency": string(addr.Currency),
				"status":   (*result.Data)["status"],
				"stuck":    (*result.Data)["stuck"],
				"amount":   (*result.Data)["received_amount"],
			}, addr.Callback)
		}

		if result.Delete {
			w.Service.Box.Delete(addr.Address)
			w.Service.Repo.Delete(addr.ID)
		}
	}
}


// Check a single address.
// Returns a Result with data or delete flag.
//  check address for: 
//    - timeout will return Delete flag + data for callback
//    - can't get needed rpc url for choosen currency or network, will return Delete flag
//    - balance check failed will return null 
//    - if balance check succeeded, but "stuck" which means that transaction is not simple as expected with Transfer event 
func (w *Worker) check(addr blockchain.Address) *Result {
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	logger.Log.Info("Checking address", slog.String("address", addr.Address))

	if addr.Timeout.Before(time.Now()) {
		logger.Log.Info("Skipping address due to timeout", slog.String("address", addr.Address))
		return &Result{
			Data: &map[string]any{
				"status":          payment.Timeout,
				"stuck":           false,
				"received_amount": "0",
			},
			Delete: true,
		}
	}

	currency, err := w.Service.Provider.GetWatcher(addr.Network, addr.Currency)
	if err != nil {
		logger.Log.Warn("No watcher for", slog.String("currency", string(addr.Currency)), slog.Any("error", err))
		return &Result{Delete: true}
	}

	isbalanced := currency.GetPendingBalance(addr.Address)
	if isbalanced {
		amount, isstuck := currency.GetLatestTx(addr.Address)
		if amount == "" {
			logger.Log.Info("No latest tx found for address", slog.String("address", addr.Address))
			return nil
		}
		return &Result{
			Data: &map[string]any{
				"status":          payment.Received,
				"stuck":           isstuck,
				"received_amount": fmt.Sprintf("%v", amount),
			},
			Delete: true,
		}
	}
	return nil
}
