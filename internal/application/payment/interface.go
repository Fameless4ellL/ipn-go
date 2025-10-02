package payment

import (
	"github.com/google/uuid"
	"go-blocker/internal/infrastructure/payment"
)

type Repository interface {
	Save(payment *payment.Payment) error
	FindByID(id uuid.UUID) (*payment.Payment, error)
	UpdateStatus(id uuid.UUID, status payment.Status, receivedAmount *string, txID *string, isContractMatch *bool) error
	ExpireWhere(predicate func(p *payment.Payment) bool) error
	ListPending() ([]payment.Payment, error)
}
