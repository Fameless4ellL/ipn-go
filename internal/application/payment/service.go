package payment

import (
	"errors"
	"fmt"

	domain "go-blocker/internal/domain/payment"
	"go-blocker/internal/infrastructure/payment"
	"go-blocker/internal/pkg/config"
	"go-blocker/internal/pkg/utils"
	"go-blocker/internal/provider"
	"go-blocker/internal/rpc"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	Repo domain.Repository
}

func NewService(repo domain.Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Create(p *WebhookRequest) (*domain.Payment, error) {
	pay, err := domain.NewPayment(p.Address, p.Currency, p.Amount, p.Timeout, p.CallbackURL)
	if err != nil {
		return nil, err
	}
	return pay, s.Repo.Save(pay)
}

func (s *Service) Status(
	id uuid.UUID,
	status domain.Status,
	receivedAmount *string,
	txID *string,
	isContractMatch *bool,
) error {
	return s.Repo.UpdateStatus(id, status, receivedAmount, txID, isContractMatch)
}

func (s *Service) ExpireTimedOutPayments() error {
	return s.Repo.ExpireWhere(func(p *domain.Payment) bool {
		return p.Status == domain.Pending && time.Now().After(p.ExpiresAt)
	})
}

func (s *Service) ListPendingPayments() ([]*domain.Payment, error) {
	return s.Repo.ListPending()
}

func (s *Service) CheckTx(req *CheckTxRequest) (*CheckTxResponse, error) {
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
		if watcher.Name() != string(req.Currency) {
			continue
		}

		amount, IsStuck := watcher.IsTransactionMatch(client, url, req.Address, req.TxID)
		if IsStuck {
			utils.Send(map[string]interface{}{
				"status":          payment.Received,
				"address":         req.Address,
				"stuck":           true,
				"received_amount": fmt.Sprintf("%v", amount),
				"txid":            fmt.Sprintf("%v", req.TxID),
				"currency":        string(watcher.Name()),
			}, url)
			return &CheckTxResponse{
				Status: domain.Received,
				Amount: amount,
			}, nil
		} else if amount != "" {
			return &CheckTxResponse{
				Status: domain.Completed,
				Amount: amount,
			}, nil
		}
	}
	return nil, errors.New("no matching transaction found")
}

func (s *Service) FindLatestTx(req *FindTxRequest) (*CheckTxResponse, error) {
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
		if watcher.Name() != string(req.Currency) {
			continue
		}

		amount, IsStuck := watcher.GetLatestTx(client, url, req.Address)
		if IsStuck {
			utils.Send(map[string]interface{}{
				"status":          payment.Received,
				"address":         req.Address,
				"stuck":           true,
				"received_amount": fmt.Sprintf("%v", amount),
				"txid":            "",
				"currency":        string(watcher.Name()),
			}, url)
			return &CheckTxResponse{
				Status: domain.Received,
				Amount: amount,
			}, nil
		} else if amount != "" {
			return &CheckTxResponse{
				Status: domain.Completed,
				Amount: amount,
			}, nil
		}
	}
	return nil, errors.New("no matching transaction found")
}
