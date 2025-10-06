package payment

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ID uuid.UUID // type safety
type Status string

const (
	Pending   Status = "pending"
	Received  Status = "received"
	Completed Status = "completed"
	Timeout   Status = "timeout"
	Failed    Status = "failed"
	Mismatch  Status = "mismatch"
)

type Payment struct {
	ID             ID
	Address        string
	Currency       string
	Amount         string
	Timeout        int
	ReceivedAmount string
	TxID           string
	CallbackURL    string
	Status         Status
	IsStuck        bool
	CreatedAt      time.Time
	ExpiresAt      time.Time
}

func NewPayment(address, currency, amount string, timeoutMinutes int, callback string) (*Payment, error) {
	if address == "" || amount == "" {
		return nil, errors.New("address and amount are required")
	}

	p := &Payment{
		ID:          ID(uuid.New()),
		Address:     address,
		Currency:    currency,
		Amount:      amount,
		Timeout:     timeoutMinutes,
		CallbackURL: callback,
		Status:      Pending,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(time.Duration(timeoutMinutes) * time.Minute),
	}
	return p, nil
}

func (p *Payment) GetID() uuid.UUID {
	return uuid.UUID(p.ID)
}
