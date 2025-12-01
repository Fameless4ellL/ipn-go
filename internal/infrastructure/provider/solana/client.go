package solana

import (
	"encoding/json"
	"fmt"
	"go-blocker/internal/infrastructure/provider/solana/types"
	"go-blocker/internal/pkg/config"
	"io"
	"net/http"
	"strings"
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
	url := fmt.Sprintf("%s%s", c.BaseURL, endpoint)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("token", config.SOLapiKey)
	resp, err := http.DefaultClient.Do(req)

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

type RPC struct {
	URL string
}

func NewRPC(url string) *RPC {
	return &RPC{
		URL: url,
	}
}

func (r *RPC) Post(method string, payload interface{}) ([]byte, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshalling payload failed: %w", err)
	}
	payloadStr := fmt.Sprintf(`{"jsonrpc":"2.0","id":1,"method":"%s","params":%s}`, method, string(body))
	resp, err := http.Post(r.URL, "application/json", strings.NewReader(payloadStr))
	if err != nil {
		return nil, fmt.Errorf("[POST]request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response failed: %w", err)
	}
	return respBody, nil
}

func (r *RPC) GetBalance(address string) (*types.RPCResponse[types.GetBalanceResponse], error) {
	resp, err := r.Post("getBalance", []interface{}{address, map[string]string{"commitment": "finalized"}})
	if err != nil {
		return nil, err
	}
	var balanceResp types.RPCResponse[types.GetBalanceResponse]
	if err := json.Unmarshal(resp, &balanceResp); err != nil {
		return nil, fmt.Errorf("[GetBalance] unmarshalling failed: %w", err)
	}
	return &balanceResp, nil
}

func (r *RPC) GetTokenAccountBalance(address string) (*types.RPCResponse[types.GetTokenAccountBalance], error) {
	resp, err := r.Post("getTokenAccountBalance", []interface{}{address, map[string]string{"commitment": "finalized"}})
	if err != nil {
		return nil, err
	}

	var balanceResp types.RPCResponse[types.GetTokenAccountBalance]
	if err := json.Unmarshal(resp, &balanceResp); err != nil {
		return nil, fmt.Errorf("[GetTokenAccountBalance] unmarshalling failed: %w", err)
	}
	return &balanceResp, nil
}

func (r *RPC) TransactionByHash(txid string) (*types.RPCResponse[types.GetTransactionResult], bool, error) {
	resp, err := r.Post("getTransaction", []interface{}{txid, map[string]interface{}{"maxSupportedTransactionVersion": 0, "encoding": "json", "commitment": "finalized"}})
	if err != nil {
		return nil, false, err
	}
	var txResp types.RPCResponse[types.GetTransactionResult]
	if err := json.Unmarshal(resp, &txResp); err != nil {
		return nil, false, fmt.Errorf("[TransactionByHash] unmarshalling failed: %w", err)
	}
	return &txResp, true, nil
}

func (r *RPC) TransactionReceipt(txid string) (*types.RPCResponse[types.GetTransactionResult], bool, error) {
	return r.TransactionByHash(txid)
}
