package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go-blocker/internal/config"
	"go-blocker/internal/utils"
)

func SendCallback(p *Payment) {
	if p.CallbackURL == "" {
		return
	}

	payload := map[string]interface{}{
		"payment_id":      p.ID,
		"status":          p.Status,
		"address":         p.Address,
		"stuck":           p.IsStuck,
		"received_amount": fmt.Sprintf("%v", p.ReceivedAmount),
		"txid":            fmt.Sprintf("%v", p.TxID),
		"currency":        p.Currency,
	}

	// make a semaphore to limit concurrent callbacks
	if p.CallbackURL != "" {
		go SendRequest(p.CallbackURL, payload)
	}
	go Telegram(payload)
}

func SendRequest(url string, payload map[string]interface{}) {
	resp, err := utils.Callback(url, payload)
	if err != nil {
		config.Log.Debugf("Failed to send callback: %v", err)
		return
	}
	defer resp.Body.Close()
}

func Telegram(payload map[string]interface{}) {
	if config.BotToken == "" || config.ChatId == "" {
		config.Log.Infoln("Telegram bot token or chat ID not set")
		return
	}

	// Format message
	message := fmt.Sprintf(
		"💸 Payment Status\nID: %v\nStatus: %v\nAddress: %v\nAmount: %v\nCurrency: %v\nTxID: %v\nStuck: %v\n",
		payload["payment_id"],
		payload["status"],
		payload["address"],
		payload["received_amount"],
		payload["currency"],
		payload["txid"],
		payload["stuck"],
	)

	// Telegram API endpoint
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.BotToken)

	// Prepare request body
	body := map[string]interface{}{
		"chat_id": config.ChatId,
		"text":    message,
	}
	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		config.Log.Infof("Telegram: Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := utils.Retry(req, 3)
	if err != nil {
		config.Log.Debugf("Telegram: Failed to send message: %v", err)
		return
	}
	defer resp.Body.Close()
}
