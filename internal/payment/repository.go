package payment

import "github.com/google/uuid"

type PaymentRepository interface {
	Save(payment *Payment) error
	FindByID(id uuid.UUID) (*Payment, error)
	UpdateStatus(id uuid.UUID, status PaymentStatus) error
	ExpireWhere(predicate func(*Payment) bool) error
}
