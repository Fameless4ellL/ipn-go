package provider

import (
	"fmt"
	"go-blocker/internal/domain/blockchain"
	"go-blocker/internal/infrastructure/provider/solana"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/pkg/utils"
	"math/big"
)

type SOL struct {
	client *solana.RPC
	scan   *solana.Client
}

func (sol *SOL) GetBalance(wallet string) *big.Int {
	balance, err := sol.client.GetBalance(wallet)
	if err != nil {
		logger.Log.Errorf("Error getting pending balance for address %s: %s", wallet, err)
		return big.NewInt(0)
	}
	return big.NewInt(int64(balance.Result.Result.Value))
}

func (sol *SOL) GetERC20Balance(abi, contract, wallet string) *big.Int {
	balance, err := sol.client.GetTokenAccountBalance(wallet)
	if err != nil {
		logger.Log.Errorf("Error getting balance for address %s: %s", wallet, err)
		return big.NewInt(0)
	}

	return big.NewInt(int64(balance.Result.Result.Value.UIAmount))
}

func (sol *SOL) TransactionByHash(txid string) (*blockchain.Transaction, error) {
	Tx, _, err := sol.client.TransactionByHash(txid)
	if err != nil {
		logger.Log.Errorf("Error getting transaction for tx %s: %s", txid, err)
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
		logger.Log.Warnf("Error getting transaction receipt for tx %s: %s", txid, err)
		return nil, err
	}
	fmt.Printf("Tx: %v\n", Tx)
	return &blockchain.Transaction{}, nil
}

func (sol *SOL) GetTx(address string) (string, error) {
	resp, err := sol.scan.GetTransactions(address)
	if err != nil {
		logger.Log.Warnf("Normal Error: %v", err)
		return "", err
	}

	return resp.Data[0].TxHash, nil
}

func (sol *SOL) GetERC20(contract, address string) (string, error) {
	resp, err := sol.scan.GetERC20(address)
	if err != nil {
		logger.Log.Warnf("Error: %v", err)
		return "", err
	}

	return resp.Data[0].TxHash, nil
}

func (evm *SOL) TraceBlock(blocknumber, address string) ([]utils.TraceResult, error) {
	results := []utils.TraceResult{}
	return results, nil
}
