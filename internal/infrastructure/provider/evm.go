package provider

import (
	"context"
	"go-blocker/internal/infrastructure/provider/etherscan"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/pkg/utils"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	ABI "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
		logger.Log.Errorf("Error getting pending balance for address %s: %s", wallet, err)
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
		log.Fatalf("CallContract %v", err)
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

func (evm *EVM) TraceBlock(Tx *types.Receipt, address string) ([]utils.TraceResult, error) {
	blocknumberInt := new(big.Int)
	blocknumberInt.SetString(Tx.BlockNumber.String(), 10)
	blocknumberHex := blocknumberInt.Text(16)
	result, err := utils.TraceBlock(evm.url, "0x"+blocknumberHex)
	if err != nil {
		logger.Log.Debugf("No healthy TraceBlock: %s", err)
		return result, err
	}
	return result, nil
}

func (evm *EVM) TransactionByHash(txid string) (*types.Transaction, error) {
	Tx, _, err := evm.client.TransactionByHash(context.Background(), common.HexToHash(txid))
	if err != nil {
		return nil, err
	}
	return Tx, nil
}

func (evm *EVM) TransactionReceipt(txid string) (*types.Receipt, error) {
	Tx, err := evm.client.TransactionReceipt(context.Background(), common.HexToHash(txid))
	if err != nil {
		logger.Log.Warnf("Error getting transaction receipt for tx %s: %s", txid, err)
		return nil, err
	}
	return Tx, nil
}

func (evm *EVM) GetTx(address string) (string, error) {
	resp, err := evm.scan.GetTransactions(address, 0, 99999999, 1, 10, "asc", etherscan.InternalTx)
	if err != nil {
		logger.Log.Warnf("Internal Error: %v", err)
		resp, err = evm.scan.GetTransactions(address, 0, 99999999, 1, 10, "asc", etherscan.NormalTx)
		if err != nil {
			logger.Log.Warnf("Normal Error: %v", err)
			return "", err
		}
		return resp.Result[0].Hash, nil
	}

	return resp.Result[0].Hash, nil
}

func (evm *EVM) GetERC20(contract, address string) (string, error) {
	resp, err := evm.scan.GetERC20(contract, address, 0, 99999999, 1, 10, "desc")
	if err != nil {
		logger.Log.Warnf("Error: %v", err)
		return "", err
	}

	return resp.Result[0].Hash, nil
}
