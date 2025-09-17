package payment

import (
	constants "go-blocker/internal/const"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Save(payment *constants.Payment) error
	FindByID(id uuid.UUID) (*constants.Payment, error)
	UpdateStatus(id uuid.UUID, status constants.PaymentStatus, receivedAmount *string, txID *string, isContractMatch *bool) error
	ExpireWhere(predicate func(p *constants.Payment) bool) error
	ListPending() ([]constants.Payment, error)
}
