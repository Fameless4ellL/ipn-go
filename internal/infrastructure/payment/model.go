package payment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Status string

const (
	Pending   Status = "pending"
	Received  Status = "received"
	Completed Status = "completed"
	Timeout   Status = "timeout"
	Failed    Status = "failed"
	Mismatch  Status = "mismatch"
)

type Payment struct {
	gorm.Model
	ID             uuid.UUID     `gorm:"type:char(36);primaryKey" json:"id"`
	Address        string        `json:"address"`
	Currency       string        `json:"currency"`
	Amount         string        `json:"amount"`
	Timeout        int           `json:"timeout"` // in minutes
	ReceivedAmount string        `json:"received_amount"`
	TxID           string        `json:"txid"`
	CallbackURL    string        `json:"callback_url"`
	Status         Status `gorm:"type:text" json:"status"`
	IsStuck        bool          `json:"Stuck"`
	CreatedAt      time.Time     `json:"created_at"`
	ExpiresAt      time.Time     `json:"expires_at"`
}

func (Payment) TableName() string {
	return "payments"
}


func (p *Payment) MakePayload() map[string]interface{} {
	payload := map[string]interface{}{
		// "payment_id":      p.ID,
		"status":          p.Status,
		"address":         p.Address,
		"stuck":           p.IsStuck,
		"received_amount": fmt.Sprintf("%v", p.ReceivedAmount),
		"txid":            fmt.Sprintf("%v", p.TxID),
		"currency":        p.Currency,
	}
	return payload
}
