package types

type TxError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tx struct {
	Slot               int           `json:"slot"`
	Fee                string        `json:"fee"`
	Status             string        `json:"status"`
	Signer             string        `json:"signer"`
	BlockTime          string        `json:"block_time"`
	TxHash             string        `json:"tx_hash"`
	ParsedInstructions []interface{} `json:"parsed_instructions"`
	ProgramIds         []string      `json:"program_ids"`
	Time               string        `json:"time"`
}

type TxListResponse struct {
	Status bool     `json:"status"`
	Data   []*Tx    `json:"data"`
	Error  *TxError `json:"errors,omitempty"`
}
