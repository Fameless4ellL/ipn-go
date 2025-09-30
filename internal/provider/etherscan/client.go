package etherscan

import (
	"fmt"
	"go-blocker/internal/config"
	"io"
	"net/http"
)

type Client struct {
	APIKey  string
	BaseURL string
}

func NewClient() *Client {
	return &Client{
		APIKey:  config.ETHapiKey,
		BaseURL: "https://api.etherscan.io/v2/api",
	}
}

func (c *Client) get(endpoint string) ([]byte, error) {
	url := fmt.Sprintf("%s%s&apikey=%s", c.BaseURL, endpoint, c.APIKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response failed: %w", err)
	}

	return body, nil
}
