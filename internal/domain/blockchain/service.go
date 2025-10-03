package domain

type BlockchainService interface {
	GetTransactionStatus(txHash string) (TransactionStatus, error)
}

type TransactionStatus string

const (
	StatusPending TransactionStatus = "pending"
	StatusSuccess TransactionStatus = "success"
	StatusFailed  TransactionStatus = "failed"
)
