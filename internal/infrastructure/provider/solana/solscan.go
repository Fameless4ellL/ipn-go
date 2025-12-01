package solana

import (
	"encoding/json"
	"fmt"
	"go-blocker/internal/infrastructure/provider/solana/types"
)

// https://pro-api.solscan.io/pro-api-docs/v2.0/reference/v2-account-transactions
func (c *Client) GetTransactions(address string) (*types.TxListResponse, error) {
	endpoint := fmt.Sprintf("?address=%s", address)
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}
	var txResp types.TxListResponse
	if err := json.Unmarshal(body, &txResp); err != nil {
		return nil, fmt.Errorf("unmarshalling failed: %w", err)
	}

	if !txResp.Status {
		return nil, fmt.Errorf("API error: %s", txResp.Error.Message)
	}

	return &txResp, nil
}

func (c *Client) GetERC20(address string) (*types.TxListResponse, error) {
	endpoint := fmt.Sprintf("?address=%s", address)
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var txResp types.TxListResponse
	if err := json.Unmarshal(body, &txResp); err != nil {
		return nil, fmt.Errorf("unmarshalling failed: %w", err)
	}

	if !txResp.Status {
		return nil, fmt.Errorf("API error: %s", txResp.Error.Message)
	}

	return &txResp, nil
}
