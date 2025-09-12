package payment_worker

import (
	"go-blocker/internal/config"
	"go-blocker/internal/payment"
	"time"
)

func Start(service *payment.PaymentService) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		
		for range ticker.C {
			err := service.ExpireTimedOutPayments()
			if err != nil {
				config.Log.Infof("[Worker] Error: %e", err)
			}
		}
	}()
}
