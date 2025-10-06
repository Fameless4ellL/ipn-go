package payment

import (
	"fmt"

	blockchain "go-blocker/internal/domain/blockchain"
	payment "go-blocker/internal/domain/payment"
	"go-blocker/internal/pkg/utils"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	Repo     payment.Repository
	Manager  blockchain.Manager
	Provider blockchain.Watcher
	Box      blockchain.Storage
}

func NewService(
	repo payment.Repository,
	m blockchain.Manager,
	p blockchain.Watcher,
	b blockchain.Storage,
) *Service {
	return &Service{Repo: repo, Manager: m, Provider: p, Box: b}
}

func (s *Service) Create(p *WebhookRequest) (*payment.Payment, error) {
	pay, err := payment.NewPayment(p.Address, p.Currency, "0", p.Timeout, p.CallbackURL)
	if err != nil {
		return nil, err
	}
	return pay, s.Repo.Save(pay)
}

func (s *Service) Status(
	id uuid.UUID,
	status payment.Status,
	receivedAmount *string,
	txID *string,
	isContractMatch *bool,
) error {
	return s.Repo.UpdateStatus(id, status, receivedAmount, txID, isContractMatch)
}

func (s *Service) ExpireTimedOutPayments() error {
	return s.Repo.ExpireWhere(func(p *payment.Payment) bool {
		return p.Status == payment.Pending && time.Now().After(p.ExpiresAt)
	})
}

func (s *Service) ListPendingPayments() ([]*payment.Payment, error) {
	return s.Repo.ListPending()
}

func (s *Service) CheckTx(req *CheckTxRequest) (*CheckTxResponse, error) {
	client, url, err := s.Manager.GetClientForChain(blockchain.Ethereum)
	if err != nil {
		return nil, err
	}

	currency, err := s.Provider.GetWatcher(blockchain.Ethereum, blockchain.CurrencyType(req.Currency))
	if err != nil {
		return nil, err
	}

	address, _ := s.Box.Get(req.Address)
	amount, IsStuck := currency.IsTransactionMatch(client, url, req.Address, req.TxID)
	if IsStuck {
		utils.Send(map[string]interface{}{
			"status":          payment.Received,
			"address":         req.Address,
			"stuck":           true,
			"received_amount": fmt.Sprintf("%v", amount),
			"txid":            fmt.Sprintf("%v", req.TxID),
			"currency":        string(currency.Name()),
		}, address.Callback)
		s.Box.Delete(req.Address)
		return &CheckTxResponse{
			Status: payment.Received,
			Amount: amount,
		}, nil
	} else if amount != "" {
		return &CheckTxResponse{
			Status: payment.Completed,
			Amount: amount,
		}, nil
	}
	return nil, fmt.Errorf("no matching transaction found for address %s and txid %s", req.Address, req.TxID)
}

func (s *Service) FindLatestTx(req *FindTxRequest) (*CheckTxResponse, error) {
	client, url, err := s.Manager.GetClientForChain(blockchain.Ethereum)
	if err != nil {
		return nil, err
	}

	currency, err := s.Provider.GetWatcher(blockchain.Ethereum, blockchain.CurrencyType(req.Currency))
	if err != nil {
		return nil, err
	}

	address, _ := s.Box.Get(req.Address)

	amount, IsStuck := currency.GetLatestTx(client, url, req.Address)
	if IsStuck {
		utils.Send(map[string]interface{}{
			"status":          payment.Received,
			"address":         req.Address,
			"stuck":           true,
			"received_amount": fmt.Sprintf("%v", amount),
			"txid":            "",
			"currency":        string(currency.Name()),
		}, address.Callback)
		s.Box.Delete(req.Address)
		return &CheckTxResponse{
			Status: payment.Received,
			Amount: amount,
		}, nil
	} else if amount != "" {
		return &CheckTxResponse{
			Status: payment.Completed,
			Amount: amount,
		}, nil
	}
	return nil, fmt.Errorf("no transactions found for address %s", req.Address)
}
