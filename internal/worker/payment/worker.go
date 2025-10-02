package payment_worker

import (
	"go-blocker/internal/application/payment"
	logger "go-blocker/internal/pkg/log"
	"time"
)

func Start(service *payment.Service) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			err := service.ExpireTimedOutPayments()
			if err != nil {
				logger.Log.Infof("[Worker] Error: %e", err)
			}
		}
	}()
}
