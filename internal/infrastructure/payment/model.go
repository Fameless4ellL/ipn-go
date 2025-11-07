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
	Network        string    `json:"network"`
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
		ID:          domain.ID(m.ID),
		Address:     m.Address,
		Network:     m.Network,
		Currency:    m.Currency,
		Amount:      m.Amount,
		Timeout:     m.Timeout,
		CallbackURL: m.CallbackURL,
		Status:      domain.Status(m.Status),
		IsStuck:     m.IsStuck,
		CreatedAt:   m.CreatedAt,
	}
}

func FromDomain(pay *domain.Payment) *PaymentModel {
	return &PaymentModel{
		ID:          uuid.UUID(pay.ID),
		Address:     pay.Address,
		Network:     pay.Network,
		Currency:    pay.Currency,
		Amount:      pay.Amount,
		Timeout:     pay.Timeout,
		CallbackURL: pay.CallbackURL,
		Status:      Status(pay.Status),
		IsStuck:     pay.IsStuck,
		CreatedAt:   pay.CreatedAt,
	}
}

func (p *PaymentModel) MakePayload() map[string]interface{} {
	payload := map[string]interface{}{
		"status":          p.Status,
		"address":         p.Address,
		"stuck":           p.IsStuck,
		"received_amount": fmt.Sprintf("%v", p.ReceivedAmount),
		"currency":        p.Currency,
		"network":         p.Network,
	}
	return payload
}
