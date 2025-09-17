package database

import (
	"go-blocker/internal/config"
	constants "go-blocker/internal/const"
	logger "go-blocker/internal/log"
	"go-blocker/internal/payment"
	"go-blocker/internal/storage"
	"math/big"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Save(p *constants.Payment) error {
	return r.db.Create(&p).Error
}

func (r *PaymentRepository) FindByID(id uuid.UUID) (*constants.Payment, error) {
	var model constants.Payment
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *PaymentRepository) UpdateStatus(
	id uuid.UUID,
	status constants.PaymentStatus,
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
	if status == constants.StatusCompleted && receivedAmount != nil {
		balance, _, err := big.ParseFloat(*receivedAmount, 10, 64, big.ToNearestEven)
		if err != nil {
			return err
		}
		if !r.isBalanceSufficient(balance, p.Amount) {
			logger.Log.Debugf("Payment %s: balance %s is less than expected %s", id, *receivedAmount, p.Amount)
			updates["status"] = string(constants.StatusMismatch)
			p.Status = constants.StatusMismatch
		}
	}

	if status != constants.StatusPending && status != constants.StatusReceived {
		storage.PaymentAddressStore.Delete(p.Address)
	}

	payment.SendCallback(p)

	return r.db.Model(&constants.Payment{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *PaymentRepository) ExpireWhere(condition func(*constants.Payment) bool) error {
	var payments []constants.Payment
	if err := r.db.Find(&payments).Error; err != nil {
		return err
	}

	for _, p := range payments {
		if condition(&p) {
			if err := r.UpdateStatus(p.ID, constants.StatusTimeout, nil, nil, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *PaymentRepository) ListPending() ([]constants.Payment, error) {
	var payments []constants.Payment
	err := r.db.Where("status = ?", "pending").Find(&payments).Error
	return payments, err
}

func (s *PaymentRepository) isBalanceSufficient(balance *big.Float, expected string) bool {
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
