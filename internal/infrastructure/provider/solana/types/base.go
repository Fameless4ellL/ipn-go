package types

type Context struct {
	Slot       int    `json:"slot"`
	ApiVersion string `json:"apiVersion"`
}

type Result[T any] struct {
	Context Context `json:"context"`
	Value   T       `json:"value"`
}

type RPCResponse[T any] struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  T      `json:"result"`
}

type RPCResponse1[T any] struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  T      `json:"result"`
}

type GetTransaction struct {
	Result Result[GetTransactionResult] `json:"result"`
}

type GetBalanceResponse struct {
	Result Result[int] `json:"result"`
}

type GetTokenAccountBalance struct {
	Result Result[GetTokenAccountBalanceValue] `json:"result"`
}
