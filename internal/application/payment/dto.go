package payment

import (
	"go-blocker/internal/domain/payment"
)

// WebhookRequest represents the incoming webhook payload
type WebhookRequest struct {
	Address string `json:"address"`
	// Network     string `json:"network"`
	Currency string `json:"currency"`
	// Amount      string `json:"amount"`
	Timeout     int    `json:"timeout"`
	CallbackURL string `json:"callback_url"`
}

// WebhookResponse represents the response to the webhook
type WebhookResponse struct {
	ID     string         `json:"id"`
	Status payment.Status `json:"status"`
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
	Status payment.Status `json:"status"`
	Amount string         `json:"amount"`
}
