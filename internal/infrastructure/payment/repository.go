package payment

import (
	"go-blocker/internal/infrastructure/notifier"
	"go-blocker/internal/pkg/config"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/storage"
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

func (r *Repository) Save(p *Payment) error {
	return r.db.Create(&p).Error
}

func (r *Repository) FindByID(id uuid.UUID) (*Payment, error) {
	var model Payment
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *Repository) UpdateStatus(
	id uuid.UUID,
	status Status,
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

	if receivedAmount != nil {
		updates["received_amount"] = *receivedAmount
		p.ReceivedAmount = *receivedAmount
	}
	if txID != nil {
		updates["TxID"] = *txID
		p.TxID = *txID
	}
	if isContractMatch != nil {
		updates["IsStuck"] = *isContractMatch
		p.IsStuck = *isContractMatch
	}

	// check isBalanceSufficient
	if status == Completed && receivedAmount != nil {
		balance, _, err := big.ParseFloat(*receivedAmount, 10, 64, big.ToNearestEven)
		if err != nil {
			return err
		}
		if !r.isBalanceSufficient(balance, p.Amount) {
			logger.Log.Debugf("Payment %s: balance %s is less than expected %s", id, *receivedAmount, p.Amount)
			updates["status"] = string(Mismatch)
			p.Status = Mismatch
		}
	}

	if status != Pending && status != Received {
		storage.PaymentAddressStore.Delete(p.Address)
	}

	notifier.Send(p.MakePayload(), p.CallbackURL)

	return r.db.Model(&Payment{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *Repository) ExpireWhere(condition func(*Payment) bool) error {
	var payments []Payment
	if err := r.db.Find(&payments).Error; err != nil {
		return err
	}

	for _, p := range payments {
		if condition(&p) {
			if err := r.UpdateStatus(p.ID, Timeout, nil, nil, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *Repository) ListPending() ([]Payment, error) {
	var payments []Payment
	err := r.db.Where("status = ?", "pending").Find(&payments).Error
	return payments, err
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
