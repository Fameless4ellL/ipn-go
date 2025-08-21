package blocker

import (
	"go-blocker/internal/payment"
	"log"
	"time"
)

func Start(service *payment.PaymentService) {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			log.Println("[Worker] Running payment expiration check...")
			err := service.ExpireTimedOutPayments()
			if err != nil {
				log.Println("[Worker] Error:", err)
			}
		}
	}()
}
