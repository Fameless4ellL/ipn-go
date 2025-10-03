package payment

import "github.com/google/uuid"

type Repository interface {
	Save(payment *Payment) error
	FindByID(id uuid.UUID) (*Payment, error)
	UpdateStatus(id uuid.UUID, status Status, receivedAmount *string, txID *string, isContractMatch *bool) error
	ExpireWhere(predicate func(p *Payment) bool) error
	ListPending() ([]*Payment, error)
}
