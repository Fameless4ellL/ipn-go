package payment

type InvalidRequest struct {
	Error string `json:"error" example:"invalid request"`
}

type InvalidAddress struct {
	Error string `json:"error" example:"Invalid address format"`
}

type FailedToFind struct {
	Error string `json:"error" example:"failed to find latest transaction"`
}
