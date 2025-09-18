package constants

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentStatus string
type PaymentStacked bool

const (
	StatusPending   PaymentStatus = "pending"
	StatusReceived  PaymentStatus = "received"
	StatusCompleted PaymentStatus = "completed"
	StatusTimeout   PaymentStatus = "timeout"
	StatusFailed    PaymentStatus = "failed"
	StatusMismatch  PaymentStatus = "mismatch"
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
	Status         PaymentStatus `gorm:"type:text" json:"status"`
	IsStuck        bool          `json:"Stuck"`
	CreatedAt      time.Time     `json:"created_at"`
	ExpiresAt      time.Time     `json:"expires_at"`
}

// WebhookRequest represents the incoming webhook payload
type WebhookRequest struct {
	Address     string `json:"address"`
	Network     string `json:"network"`
	Currency    string `json:"currency"`
	Amount      string `json:"amount"`
	Timeout     int    `json:"timeout"`
	CallbackURL string `json:"callback_url"`
}

// WebhookResponse represents the response to the webhook
type WebhookResponse struct {
	ID     string        `json:"id"`
	Status PaymentStatus `json:"status"`
}

type CurrencyType string

const (
	ETH  CurrencyType = "ETH"
	USDT CurrencyType = "USDT"
	USDC CurrencyType = "USDC"
)

// CheckTxRequest represents the request to check a transaction
type CheckTxRequest struct {
	Address  string       `json:"address" example:"0xabc123..."`
	Currency CurrencyType `json:"currency" example:"USDT" description:"Token symbol (e.g., ETH, USDC, USDT)"`
	TxID     string       `json:"txid" example:"0xabc123..." description:"The transaction ID to check"`
}

type FindTxRequest struct {
	Address  string       `json:"address" example:"0xabc123..." format:"hex"`
	Currency CurrencyType `json:"currency" example:"USDT" description:"Token symbol (e.g., ETH, USDC, USDT)"`
}

// CheckTxResponse represents the response from checking a transaction
type CheckTxResponse struct {
	Status PaymentStatus `json:"status"`
	Amount string        `json:"amount"`
}
