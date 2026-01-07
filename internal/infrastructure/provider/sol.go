package provider

import (
	"fmt"
	"go-blocker/internal/domain/blockchain"
	"go-blocker/internal/infrastructure/provider/solana"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/pkg/utils"
	"log/slog"
	"math/big"
)

type SOL struct {
	client *solana.RPC
	scan   *solana.Client
}

func (sol *SOL) GetBalance(wallet string) *big.Int {
	balance, err := sol.client.GetBalance(wallet)
	if err != nil {
		logger.Log.Error(
			"Error getting pending balance",
			slog.Group("SOL",
				slog.String("address", wallet),
				slog.Any("error", err),
			),
		)
		return big.NewInt(0)
	}
	return big.NewInt(int64(balance.Result.Result.Value))
}

func (sol *SOL) GetERC20Balance(abi, contract, wallet string) *big.Int {
	balance, err := sol.client.GetTokenAccountBalance(wallet)
	if err != nil {
		logger.Log.Error(
			"Error getting balance",
			slog.Group("SOL",
				slog.String("address", wallet),
				slog.Any("error", err),
			),
		)
		return big.NewInt(0)
	}

	return big.NewInt(int64(balance.Result.Result.Value.UIAmount))
}

func (sol *SOL) TransactionByHash(txid string) (*blockchain.Transaction, error) {
	Tx, _, err := sol.client.TransactionByHash(txid)
	if err != nil {
		logger.Log.Error(
			"Error getting transaction",
			slog.Group("SOL",
				slog.String("txid", txid),
				slog.Any("error", err),
			),
		)
		return nil, err
	}
	fmt.Printf("Tx: %v\n", Tx)
	return &blockchain.Transaction{
		BlockNumber: big.NewInt(0),
	}, nil
}

func (sol *SOL) TransactionReceipt(txid string) (*blockchain.Transaction, error) {
	Tx, _, err := sol.client.TransactionReceipt(txid)
	if err != nil {
		logger.Log.Error(
			"Error getting transaction receipt",
			slog.Group("SOL",
				slog.String("txid", txid),
				slog.Any("error", err),
			),
		)
		return nil, err
	}
	fmt.Printf("Tx: %v\n", Tx)
	return &blockchain.Transaction{}, nil
}

func (sol *SOL) GetTx(address string) (string, error) {
	resp, err := sol.scan.GetTransactions(address)
	if err != nil {
		logger.Log.Error(
			"Error getting transactions",
			slog.Group("SOL",
				slog.String("address", address),
				slog.Any("error", err),
			),
		)
		return "", err
	}

	return resp.Data[0].TxHash, nil
}

func (sol *SOL) GetERC20(contract, address string) (string, error) {
	resp, err := sol.scan.GetERC20(address)
	if err != nil {
		logger.Log.Error(
			"Error getting ERC20 transactions",
			slog.Group("SOL",
				slog.String("contract", contract),
				slog.String("address", address),
				slog.Any("error", err),
			),
		)
		return "", err
	}

	return resp.Data[0].TxHash, nil
}

func (evm *SOL) TraceBlock(blocknumber, address string) ([]utils.TraceResult, error) {
	results := []utils.TraceResult{}
	return results, nil
}
