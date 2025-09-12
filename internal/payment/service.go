package payment

import (
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Create(p *WebhookRequest) (*Payment, error) {
	payment := &Payment{
		ID:          uuid.New(),
		Address:     p.Address,
		Currency:    p.Currency,
		Amount:      p.Amount,
		Timeout:     p.Timeout,
		CallbackURL: p.CallbackURL,
		Status:      StatusPending,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(time.Duration(p.Timeout) * time.Minute),
	}
	return payment, s.repo.Save(payment)
}

func (s *PaymentService) Status(
	id uuid.UUID,
	status PaymentStatus,
	receivedAmount *string,
	txID *string,
	isContractMatch *bool,
) error {
	return s.repo.UpdateStatus(id, status, receivedAmount, txID, isContractMatch)
}

func (s *PaymentService) ExpireTimedOutPayments() error {
	return s.repo.ExpireWhere(func(p *Payment) bool {
		return p.Status == StatusPending && time.Now().After(p.ExpiresAt)
	})
}

func (s *PaymentService) ListPendingPayments() ([]Payment, error) {
	return s.repo.ListPending()
}
