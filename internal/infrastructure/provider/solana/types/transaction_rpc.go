package types

type AccountKey struct {
	Pubkey   string `json:"pubkey"`
	Signer   bool   `json:"signer"`
	Source   string `json:"source"`
	Writable bool   `json:"writable"`
}

type Info struct {
	Destination string `json:"destination"`
	Lamports    int    `json:"lamports"`
	Source      string `json:"source"`
}

type Parsed struct {
	Info Info   `json:"info"`
	Type string `json:"type"`
}

type Instruction struct {
	Parsed      Parsed `json:"parsed"`
	Program     string `json:"program"`
	ProgramId   string `json:"programId"`
	StackHeight string `json:"stackHeight"`
}

type Message struct {
	AccountKeys         []any  `json:"accountKeys"`
	AddressTableLookups []any  `json:"addressTableLookups"`
	Header              any    `json:"header"`
	Instructions        []any  `json:"instructions"`
	RecentBlockhash     string `json:"recentBlockhash"`
}

type Transaction struct {
	Signatures []string `json:"signatures"`
	Message    Message  `json:"message"`
}

type Status struct {
	Ok *string `json:"Ok"`
}

type Meta struct {
	any
	// ComputeUnitsConsumed int           `json:"computeUnitsConsumed"`
	// CostUnits            int           `json:"costUnits"`
	// Err                  *string       `json:"err"`
	// Fee                  int           `json:"fee"`
	// InnerInstructions    []interface{} `json:"innerInstructions"`
	// LoadedAddresses      []any         `json:"loadedAddresses"`
	LogMessages       []string      `json:"logMessages"`
	PostBalances      []int         `json:"postBalances"`
	PostTokenBalances []interface{} `json:"postTokenBalances"`
	PreBalances       []int         `json:"preBalances"`
	PreTokenBalances  []*any        `json:"preTokenBalances"`
	// Rewards              []*any        `json:"rewards"`
	// Status               Status        `json:"status"`
}

type GetTransactionResult struct {
	Blocktime   int  `json:"blocktime"`
	Meta        Meta `json:"meta"`
	Slot        int  `json:"slot"`
	Transaction any  `json:"transaction"`
	Version     any  `json:"version"`
}
