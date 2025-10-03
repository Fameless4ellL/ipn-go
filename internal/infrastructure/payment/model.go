package payment

import (
	"fmt"
	domain "go-blocker/internal/domain/payment"
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

type PaymentModel struct {
	gorm.Model
	ID             uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	Address        string    `json:"address"`
	Currency       string    `json:"currency"`
	Amount         string    `json:"amount"`
	Timeout        int       `json:"timeout"` // in minutes
	ReceivedAmount string    `json:"received_amount"`
	TxID           string    `json:"txid"`
	CallbackURL    string    `json:"callback_url"`
	Status         Status    `gorm:"type:text" json:"status"`
	IsStuck        bool      `json:"Stuck"`
	CreatedAt      time.Time `json:"created_at"`
	ExpiresAt      time.Time `json:"expires_at"`
}

func (PaymentModel) TableName() string {
	return "payments"
}

func (m *PaymentModel) ToDomain() *domain.Payment {
	return &domain.Payment{
		ID:             domain.ID(m.ID),
		Address:        m.Address,
		Currency:       m.Currency,
		Amount:         m.Amount,
		Timeout:        m.Timeout,
		ReceivedAmount: m.ReceivedAmount,
		TxID:           m.TxID,
		CallbackURL:    m.CallbackURL,
		Status:         domain.Status(m.Status),
		IsStuck:        m.IsStuck,
		CreatedAt:      m.CreatedAt,
		ExpiresAt:      m.ExpiresAt,
	}
}

func FromDomain(pay *domain.Payment) *PaymentModel {
	return &PaymentModel{
		ID:             uuid.UUID(pay.ID),
		Address:        pay.Address,
		Currency:       pay.Currency,
		Amount:         pay.Amount,
		Timeout:        pay.Timeout,
		ReceivedAmount: pay.ReceivedAmount,
		TxID:           pay.TxID,
		CallbackURL:    pay.CallbackURL,
		Status:         Status(pay.Status),
		IsStuck:        pay.IsStuck,
		CreatedAt:      pay.CreatedAt,
		ExpiresAt:      pay.ExpiresAt,
	}
}

func (p *PaymentModel) MakePayload() map[string]interface{} {
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
