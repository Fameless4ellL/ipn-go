package handlers

import (
	"fmt"
	application "go-blocker/internal/application/payment"
	blockchain "go-blocker/internal/domain/blockchain"

	tele "gopkg.in/telebot.v4"
)

type TGHandler struct {
	Service *application.Service
}

func NewTGHandler(s *application.Service) *TGHandler {
	return &TGHandler{Service: s}
}

func (tg TGHandler) CheckTx(c tele.Context) error {
	args := c.Args()
	if len(args) < 4 {
		return c.Send("Usage: /check <network> <currency> <address> <txid>")
	}
	network, currency, address, txid := args[0], args[1], args[2], args[3]

	req := &application.CheckTxRequest{
		Network:  network,
		Currency: blockchain.CurrencyType(currency),
		Address:  address,
		TxID:     txid,
	}
	resp, err := tg.Service.CheckTx(req)
	if err != nil {
		return c.Send("Error checking transaction: " + err.Error())
	}
	return c.Send(fmt.Sprintf("Status: %s, Amount: %s", resp.Status, resp.Amount))
}

func (tg TGHandler) FindTx(c tele.Context) error {
	args := c.Args()
	if len(args) < 3 {
		return c.Send("Usage: /find <network> <currency> <address>")
	}
	network, currency, address := args[0], args[1], args[2]

	req := &application.FindTxRequest{
		Network:  network,
		Currency: blockchain.CurrencyType(currency),
		Address:  address,
	}
	resp, err := tg.Service.FindLatestTx(req)
	if err != nil {
		return c.Send("Error finding latest transaction: " + err.Error())
	}
	return c.Send(fmt.Sprintf("Status: %s, Amount: %s", resp.Status, resp.Amount))
}
