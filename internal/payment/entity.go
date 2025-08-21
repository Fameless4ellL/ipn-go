package payment

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string

const (
	StatusPending   PaymentStatus = "pending"
	StatusReceived  PaymentStatus = "received"
	StatusCompleted PaymentStatus = "completed"
	StatusTimeout   PaymentStatus = "timeout"
)

type Payment struct {
	gorm.Model
	ID          uuid.UUID     `gorm:"type:char(36);primaryKey" json:"id"`
	Address     string        `json:"address"`
	Currency    string        `json:"currency"`
	Amount      string        `json:"amount"`
	Timeout     int           `json:"timeout"` // in minutes
	CallbackURL string        `json:"callback_url"`
	Status      PaymentStatus `gorm:"type:text" json:"status"`
	CreatedAt   time.Time     `json:"created_at"`
	ExpiresAt   time.Time     `json:"expires_at"`
}
