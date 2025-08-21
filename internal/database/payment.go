package database

import (
	"go-blocker/internal/payment"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Save(p *payment.Payment) error {
	return r.db.Create(&p).Error
}

func (r *PaymentRepository) FindByID(id uuid.UUID) (*payment.Payment, error) {
	var model payment.Payment
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *PaymentRepository) UpdateStatus(id uuid.UUID, status payment.PaymentStatus) error {
	return r.db.Model(&payment.Payment{}).
		Where("id = ?", id).
		Update("status", string(status)).Error
}

func (r *PaymentRepository) ExpireWhere(condition func(*payment.Payment) bool) error {
	var payments []payment.Payment
	if err := r.db.Find(&payments).Error; err != nil {
		return err
	}

	for _, p := range payments {
		if condition(&p) {
			p.Status = payment.StatusTimeout
			if err := r.UpdateStatus(p.ID, p.Status); err != nil {
				return err
			}
		}
	}
	return nil
}
