package utils

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"go-blocker/internal/pkg/config"
	logger "go-blocker/internal/pkg/log"

	"resty.dev/v3"
)

// police for 200 to 300 is true otherwise false
func CircuitBreaker3xxPolicy(resp *http.Response) bool {
	return resp.StatusCode > 299
}

func NewClient() *resty.Client {
	var cb = resty.NewCircuitBreaker().
		SetPolicies(CircuitBreaker3xxPolicy)

	return resty.New().
		SetCircuitBreaker(cb).
		SetTimeout(15 * time.Second)
}

func Send(payload map[string]any, url string) {
	if url != "" {
		go SendRequest(url, payload)
	}
	if payload["stuck"] == true {
		go Telegram(payload, config.ChatId)
	}
}

func SendRequest(url string, payload map[string]any) {
	resp, err := Callback(url, payload)
	if err != nil {
		logger.Log.Debug("failed to send callback", slog.Any("error", err))
		return
	}
	defer resp.Body.Close()
}

func Telegram(payload map[string]any, chatid string) {
	if config.BotToken == "" || config.ChatId == "" {
		logger.Log.Info("Telegram bot token or chat ID not set")
		return
	}

	// Format message
	message := fmt.Sprintf(
		"💸 Payment Status\nStatus: %v\nAddress: %v\nAmount: %v\nCurrency: %v\n",
		payload["status"],
		payload["address"],
		payload["received_amount"],
		payload["currency"],
	)

	if network, ok := payload["network"]; ok {
		message += fmt.Sprintf("Network: %v\n", network)
	}

	if _, ok := payload["stuck"]; ok && payload["stuck"] == true {
		message += "⚠️ This payment was marked as stuck.\n"
	}

	// Telegram API endpoint
	url := fmt.Sprintf("%s/bot%s/sendMessage", config.TG_BASE_URL, config.BotToken)

	Client := NewClient()
	req, err := Client.R().
		SetBody(map[string]any{
			"chat_id": chatid,
			"text":    message,
		}).
		SetHeader("authorization", config.TELEGRAM_AUTH_BASE_URL).
		SetHeader("Content-Type", "application/json").
		Post(url)
	if err != nil {
		logger.Log.Debug("failed to send telegram message", slog.Any("error", err))
		return
	}
	defer req.Body.Close()
}

func Callback(url string, payload map[string]any) (*resty.Response, error) {
	Client := NewClient()
	req, err := Client.R().
		SetBody(payload).
		SetHeader("Content-Type", "application/json").
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	return req, nil
}
