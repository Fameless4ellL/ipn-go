package etherscan

import (
	"encoding/json"
	"fmt"
	"strings"
)

func (c *Client) GetTransactions(address string, startBlock, endBlock int, page int, offset int, sort string) (*TxListResponse, error) {
	endpoint := fmt.Sprintf("?chainid=1&module=account&action=txlist&address=%s&startblock=%d&endblock=%d&page=%d&offset=%d&sort=%s", address, startBlock, endBlock, page, offset, sort)
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var txResp TxListResponse
	if err := json.Unmarshal(body, &txResp); err != nil {
		return nil, fmt.Errorf("unmarshalling failed: %w", err)
	}

	if txResp.Status != "1" {
		return nil, fmt.Errorf("API error: %s", txResp.Message)
	}

	return &txResp, nil
}

func (c *Client) GetERC20(contractaddress string, address string, startBlock, endBlock int, page int, offset int, sort string) (*TxListResponse, error) {
	endpoint := fmt.Sprintf("?chainid=1&module=account&action=tokentx&contractaddress=%s&address=%s&startblock=%d&endblock=%d&page=%d&offset=%d&sort=%s", strings.ToLower(contractaddress), strings.ToLower(address), startBlock, endBlock, page, offset, sort)
	body, err := c.get(endpoint)
	if err != nil {
		return nil, err
	}

	var txResp TxListResponse
	if err := json.Unmarshal(body, &txResp); err != nil {
		return nil, fmt.Errorf("unmarshalling failed: %w", err)
	}

	if txResp.Status != "1" {
		return nil, fmt.Errorf("API error: %s", txResp.Message)
	}

	return &txResp, nil
}
