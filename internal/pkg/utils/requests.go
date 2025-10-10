package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	"go-blocker/internal/pkg/config"
	logger "go-blocker/internal/pkg/log"
)

var Client = &http.Client{
	Timeout: 5 * time.Second,

	Transport: &http.Transport{
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
	},
}

func Send(payload map[string]interface{}, url string) {
	// make a semaphore to limit concurrent callbacks
	if url != "" {
		go SendRequest(url, payload)
	}
	if payload["stuck"] == true {
		go Telegram(payload, config.ChatId)
	}
}

func SendRequest(url string, payload map[string]interface{}) {
	resp, err := Callback(url, payload)
	if err != nil {
		logger.Log.Debugf("Failed to send callback: %v", err)
		return
	}
	defer resp.Body.Close()
}

func Telegram(payload map[string]interface{}, chatid string) {
	if config.BotToken == "" || config.ChatId == "" {
		logger.Log.Infoln("Telegram bot token or chat ID not set")
		return
	}

	// Format message
	message := fmt.Sprintf(
		"💸 Payment Status\nStatus: %v\nAddress: %v\nAmount: %v\nCurrency: %v\nStuck: %v\n",
		payload["status"],
		payload["address"],
		payload["received_amount"],
		payload["currency"],
		payload["stuck"],
	)

	// Telegram API endpoint
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", config.BotToken)

	// Prepare request body
	body := map[string]interface{}{
		"chat_id": chatid,
		"text":    message,
	}
	jsonData, _ := json.Marshal(body)

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		logger.Log.Infof("Telegram: Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := Retry(req, 3)
	if err != nil {
		logger.Log.Debugf("Telegram: Failed to send message: %v", err)
		return
	}
	defer resp.Body.Close()
}

func Callback(url string, payload map[string]interface{}) (*http.Response, error) {
	jsonData, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := Retry(req, 3)

	if err != nil {
		return nil, fmt.Errorf("failed to retry request: %v", err)
	}
	return resp, nil
}

func Retry(req *http.Request, maxRetries int) (*http.Response, error) {
	for i := 0; i < maxRetries; i++ {
		resp, err := Client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			return resp, nil
		}
		time.Sleep(time.Duration(i+1) * time.Second) // backoff
	}
	return nil, fmt.Errorf("request failed after %d retries", maxRetries)
}
