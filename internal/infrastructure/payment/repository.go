package payment

import (
	domain "go-blocker/internal/domain/payment"
	"go-blocker/internal/infrastructure/storage"
	"go-blocker/internal/pkg/config"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/pkg/utils"
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(p *domain.Payment) error {
	model := FromDomain(p)
	return r.db.Create(&model).Error
}

func (r *Repository) FindByID(id uuid.UUID) (*domain.Payment, error) {
	var model PaymentModel
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return model.ToDomain(), nil
}

func (r *Repository) UpdateStatus(
	id uuid.UUID,
	status domain.Status,
	receivedAmount *string,
	txID *string,
	isContractMatch *bool,
) error {
	p, err := r.FindByID(id)
	if err != nil {
		return err
	}
	updates := map[string]interface{}{
		"status": string(status),
	}
	p.Status = status
	model := FromDomain(p)

	if isContractMatch != nil {
		updates["IsStuck"] = *isContractMatch
		p.IsStuck = *isContractMatch
	}

	// check isBalanceSufficient
	if model.Status == Completed && receivedAmount != nil {
		balance, _, err := big.ParseFloat(*receivedAmount, 10, 64, big.ToNearestEven)
		if err != nil {
			return err
		}
		if !r.isBalanceSufficient(balance, p.Amount) {
			logger.Log.Debugf("Payment %s: balance %s is less than expected %s", id, *receivedAmount, p.Amount)
			updates["status"] = string(Mismatch)
			model.Status = Mismatch
		}
	}

	utils.Send(model.MakePayload(), p.CallbackURL)

	return r.db.Model(&PaymentModel{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *Repository) ExpireWhere(condition func(*domain.Payment) bool) error {
	var payments []PaymentModel
	if err := r.db.Find(&payments).Error; err != nil {
		return err
	}

	for _, p := range payments {
		if condition(p.ToDomain()) {
			if err := r.UpdateStatus(p.ID, domain.Timeout, nil, nil, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Repository) List() ([]*domain.Payment, error) {
	var models []PaymentModel
	err := r.db.Find(&models).Error
	domainPayments := make([]*domain.Payment, len(models))
	for i, model := range models {
		domainPayments[i] = model.ToDomain()
	}
	return domainPayments, err
}

func (r *Repository) Delete(id uuid.UUID) {
	r.db.Delete(&PaymentModel{}, "id = ?", id)
}

func (s *Repository) isBalanceSufficient(balance *big.Float, expected string) bool {
	expectedBig, err := new(big.Float).SetString(expected)
	if !err {
		logger.Log.Errorf("Invalid expected amount format: %s", expected)
		return false
	}
	tolerance := new(big.Float).SetFloat64(config.BalanceTolerance)

	// Минимально требуемый баланс = ожидаемая сумма - погрешность
	minRequired := new(big.Float).Sub(expectedBig, tolerance)

	return balance.Cmp(minRequired) >= 0
}
