package payment

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
