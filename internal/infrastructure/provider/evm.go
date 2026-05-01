package provider

import (
	"context"
	"go-blocker/internal/domain/blockchain"
	"go-blocker/internal/infrastructure/provider/etherscan"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/pkg/utils"
	"log"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	ABI "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EVM struct {
	client *ethclient.Client
	scan   *etherscan.Client
	url    string
}

func (evm *EVM) GetBalance(wallet string) *big.Int {
	balance, err := evm.client.PendingBalanceAt(context.Background(), common.HexToAddress(wallet))
	if err != nil {
		logger.Log.Error(
			"Error getting pending balance for address",
			slog.String("wallet", wallet),
			slog.Any("error", err),
		)
		return big.NewInt(0)
	}
	return balance
}

func (evm *EVM) GetERC20Balance(abi, contract, wallet string) *big.Int {
	address := common.HexToAddress(wallet)
	contractAddr := common.HexToAddress(contract)

	parsedABI, err := ABI.JSON(strings.NewReader(abi))
	if err != nil {
		log.Fatalf("JSON %v", err)
	}

	data, err := parsedABI.Pack("balanceOf", address)
	if err != nil {
		log.Fatalf("Pack %v", err)
	}

	msg := ethereum.CallMsg{
		To:   &contractAddr,
		Data: data,
	}

	result, err := evm.client.CallContract(context.Background(), msg, nil)
	if err != nil {
		logger.Log.Warn("CallContract %v", slog.Any("error", err))
		return big.NewInt(0)
	}

	// Unpack returns []interface{} for output parameters
	out, err := parsedABI.Unpack("balanceOf", result)
	if err == nil && len(out) > 0 {
		switch v := out[0].(type) {
		case *big.Int:
			return v
		case big.Int:
			b := new(big.Int)
			b.Set(&v)
			return b
		}
	}

	// Fallback: try UnpackIntoInterface
	var balance *big.Int
	err = parsedABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		log.Fatalf("%v", err)
	}
	return balance
}

func (evm *EVM) TraceBlock(blocknumber, address string) ([]utils.TraceResult, error) {
	result, err := utils.TraceBlock(evm.url, "0x"+blocknumber)
	if err != nil {
		logger.Log.Error("No healthy TraceBlock", slog.Any("error", err))
		return result, err
	}
	return result, nil
}

func (evm *EVM) TransactionByHash(txid string) (*blockchain.Transaction, error) {
	Tx, _, err := evm.client.TransactionByHash(context.Background(), common.HexToHash(txid))
	if err != nil {
		return nil, err
	}

	return &blockchain.Transaction{
		BlockNumber:     nil,
		ContractAddress: Tx.To().String(),
		Hash:            Tx.Hash().Hex(),
		Logs:            []*blockchain.Logs{},
		Value:           Tx.Value(),
	}, nil
}

func (evm *EVM) TransactionReceipt(txid string) (*blockchain.Transaction, error) {
	Tx, err := evm.client.TransactionReceipt(context.Background(), common.HexToHash(txid))
	if err != nil {
		logger.Log.Error("Error getting transaction receipt", slog.String("txid", txid), slog.Any("error", err))
		return nil, err
	}

	logs := make([]*blockchain.Logs, 0, len(Tx.Logs))
	for _, log := range Tx.Logs {
		logs = append(logs, &blockchain.Logs{

			Address: log.Address.Hex(),
			Topics:  log.Topics,
			Data:    log.Data,
		})
	}

	return &blockchain.Transaction{
		BlockNumber:     Tx.BlockNumber,
		ContractAddress: Tx.ContractAddress.String(),
		Hash:            Tx.TxHash.Hex(),
		Logs:            logs,
		Value:           big.NewInt(0),
	}, nil
}

func (evm *EVM) GetTx(address string) (string, error) {
	resp, err := evm.scan.GetTransactions(address, 0, 99999999, 1, 10, "asc", etherscan.InternalTx)
	if err != nil {
		logger.Log.Error("Internal Error", slog.Any("error", err))
		resp, err = evm.scan.GetTransactions(address, 0, 99999999, 1, 10, "asc", etherscan.NormalTx)
		if err != nil {
			logger.Log.Error("Normal Error", slog.Any("error", err))
			return "", err
		}
		return resp.Result[0].Hash, nil
	}

	return resp.Result[0].Hash, nil
}

func (evm *EVM) GetERC20(contract, address string) (string, error) {
	resp, err := evm.scan.GetERC20(contract, address, 0, 99999999, 1, 10, "desc")
	if err != nil {
		logger.Log.Error("Error", slog.Any("error", err))
		return "", err
	}

	return resp.Result[0].Hash, nil
}
