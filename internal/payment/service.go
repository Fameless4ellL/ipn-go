package payment

import (
	"errors"
	"go-blocker/internal/config"
	constants "go-blocker/internal/const"
	"go-blocker/internal/provider"
	"go-blocker/internal/rpc"
	"time"

	"github.com/google/uuid"
)

type PaymentService struct {
	repo PaymentRepository
}

func NewPaymentService(repo PaymentRepository) *PaymentService {
	return &PaymentService{repo: repo}
}

func (s *PaymentService) Create(p *constants.WebhookRequest) (*constants.Payment, error) {
	payment := &constants.Payment{
		ID:          uuid.New(),
		Address:     p.Address,
		Currency:    p.Currency,
		Amount:      p.Amount,
		Timeout:     p.Timeout,
		CallbackURL: p.CallbackURL,
		Status:      constants.StatusPending,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(time.Duration(p.Timeout) * time.Minute),
	}
	return payment, s.repo.Save(payment)
}

func (s *PaymentService) Status(
	id uuid.UUID,
	status constants.PaymentStatus,
	receivedAmount *string,
	txID *string,
	isContractMatch *bool,
) error {
	return s.repo.UpdateStatus(id, status, receivedAmount, txID, isContractMatch)
}

func (s *PaymentService) ExpireTimedOutPayments() error {
	return s.repo.ExpireWhere(func(p *constants.Payment) bool {
		return p.Status == constants.StatusPending && time.Now().After(p.ExpiresAt)
	})
}

func (s *PaymentService) ListPendingPayments() ([]constants.Payment, error) {
	return s.repo.ListPending()
}

func (s *PaymentService) CheckTx(req *constants.CheckTxRequest) (*constants.CheckTxResponse, error) {
	manager := rpc.NewManager(config.Nodes)
	client, url, err := manager.GetClientForChain(rpc.Ethereum)
	if err != nil {
		return nil, err
	}

	group := map[rpc.ChainType][]provider.CurrencyWatcher{
		rpc.Ethereum: {
			&provider.ETH{},
			&provider.USDT{},
			&provider.USDC{},
		},
	}[rpc.Ethereum]

	for _, watcher := range group {
		if watcher.Name() != req.Currency {
			continue
		}

		amount, IsStuck := watcher.IsTransactionMatch(client, url, req)
		if IsStuck {
			SendCallback(&constants.Payment{
				ID:             uuid.Nil,
				Address:        req.Address,
				Currency:       string(watcher.Name()),
				Amount:         amount,
				Timeout:        0,
				CallbackURL:    "",
				Status:         constants.StatusCompleted,
				CreatedAt:      time.Time{},
				ExpiresAt:      time.Time{},
				ReceivedAmount: amount,
				TxID:           req.TxID,
				IsStuck:        true,
			})
			return &constants.CheckTxResponse{
				Status: constants.StatusReceived,
				Amount: amount,
			}, nil
		} else if amount != "" {
			return &constants.CheckTxResponse{
				Status: constants.StatusCompleted,
				Amount: amount,
			}, nil
		}
	}
	return nil, errors.New("no matching transaction found")
}

func (s *PaymentService) FindLatestTx(req *constants.FindTxRequest) (*constants.CheckTxResponse, error) {
	manager := rpc.NewManager(config.Nodes)
	client, url, err := manager.GetClientForChain(rpc.Ethereum)
	if err != nil {
		return nil, err
	}
	group := map[rpc.ChainType][]provider.CurrencyWatcher{
		rpc.Ethereum: {
			&provider.ETH{},
			&provider.USDT{},
			&provider.USDC{},
		},
	}[rpc.Ethereum]
	for _, watcher := range group {
		if watcher.Name() != req.Currency {
			continue
		}

		amount, IsStuck := watcher.GetLatestTx(client, url, *req)
		if IsStuck {
			SendCallback(&constants.Payment{
				ID:             uuid.Nil,
				Address:        req.Address,
				Currency:       string(watcher.Name()),
				Amount:         amount,
				Timeout:        0,
				CallbackURL:    "",
				Status:         constants.StatusCompleted,
				CreatedAt:      time.Time{},
				ExpiresAt:      time.Time{},
				ReceivedAmount: amount,
				TxID:           "",
				IsStuck:        true,
			})
			return &constants.CheckTxResponse{
				Status: constants.StatusReceived,
				Amount: amount,
			}, nil
		} else if amount != "" {
			return &constants.CheckTxResponse{
				Status: constants.StatusCompleted,
				Amount: amount,
			}, nil
		}
	}
	return nil, errors.New("no matching transaction found")
}
