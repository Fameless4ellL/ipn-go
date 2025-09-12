package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
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
