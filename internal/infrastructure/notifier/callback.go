package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"go-blocker/internal/pkg/config"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/utils"
)

func Send(payload map[string]interface{}, url string) {
	// make a semaphore to limit concurrent callbacks
	if url != "" {
		go SendRequest(url, payload)
	}
	go Telegram(payload)
}

func SendRequest(url string, payload map[string]interface{}) {
	resp, err := utils.Callback(url, payload)
	if err != nil {
		logger.Log.Debugf("Failed to send callback: %v", err)
		return
	}
	defer resp.Body.Close()
}

func Telegram(payload map[string]interface{}) {
	if config.BotToken == "" || config.ChatId == "" {
		logger.Log.Infoln("Telegram bot token or chat ID not set")
		return
	}

	// Format message
	message := fmt.Sprintf(
		"💸 Payment Status\nStatus: %v\nAddress: %v\nAmount: %v\nCurrency: %v\nTxID: %v\nStuck: %v\n",
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
		logger.Log.Infof("Telegram: Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := utils.Retry(req, 3)
	if err != nil {
		logger.Log.Debugf("Telegram: Failed to send message: %v", err)
		return
	}
	defer resp.Body.Close()
}
