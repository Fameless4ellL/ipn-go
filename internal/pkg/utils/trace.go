package utils

import (
	"encoding/json"
	logger "go-blocker/internal/pkg/log"
	"io"
)

type Action struct {
	From     string `json:"from"`
	CallType string `json:"callType"`
	Gas      string `json:"gas"`
	Input    string `json:"input"`
	To       string `json:"to"`
	Value    string `json:"value"`
}

// TraceResult represents the structure of the trace response
type TraceResult struct {
	Action              Action `json:"action"`
	Subtraces           int    `json:"subtraces"`
	TraceAddress        []int  `json:"traceAddress"`
	TransactionHash     string `json:"transactionHash"`
	TransactionPosition int    `json:"transactionPosition"`
	Type                string `json:"type"`
}

// JSONRPCResponse wraps the full response
type JSONRPCResponse struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Result  []TraceResult `json:"result"`
}

func TraceBlock(node string, blocknum string) ([]TraceResult, error) {
	// Prepare the request payload
	payload := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "trace_block",
		"params":  []interface{}{blocknum},
		"id":      1,
	}
	resp, err := Callback(node, payload)
	if err != nil {
		logger.Log.Debugf("Failed to send callback: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Log.Debugf("Failed to read callback: %v", err)
		return nil, err
	}

	var rpcResp JSONRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		logger.Log.Debugf("Failed to parse body callback: %v", err)
		return nil, err
	}

	return rpcResp.Result, nil
}
