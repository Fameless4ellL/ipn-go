package provider

import (
	"go-blocker/internal/domain/blockchain"
	logger "go-blocker/internal/pkg/log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type ERC20[T blockchain.API] struct {
	client   T
	Name     blockchain.CurrencyType
	Address  common.Address
	Decimals int
}

func (w *ERC20[T]) GetName() blockchain.CurrencyType {
	return w.Name
}

func (w *ERC20[T]) Abi() string {
	abiJSON := `[{"constant":true,"inputs":[{"name":"who","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"}]`
	return abiJSON
}

func (w *ERC20[T]) GetPendingBalance(wallet string) bool {
	balance := w.client.GetERC20Balance(w.Abi(), w.Address.String(), wallet)
	Balance := new(big.Float).Quo(
		new(big.Float).SetInt(balance),
		big.NewFloat(float64(w.Decimals)),
	)
	return Balance.Cmp(big.NewFloat(0)) > 0
}

func (w *ERC20[T]) Checklogs(tx *types.Receipt, address string) (string, bool) {
	for _, log := range tx.Logs {

		if len(log.Topics) < 3 {
			continue
		}

		toTopic := log.Topics[2]
		if len(toTopic.Bytes()) != 32 {
			continue
		}

		to := common.BytesToAddress(toTopic.Bytes()[12:])
		if to != common.HexToAddress(address) {
			logger.Log.Debugf("[%s]: to address not match, expected %s, got %s", w.Name, address, to.Hex())
			continue
		}

		// Convert decimals to big.Float: 10^decimals
		scale := new(big.Float).SetFloat64(1)
		for i := 0; i < w.Decimals; i++ {
			scale.Mul(scale, big.NewFloat(10))
		}

		value := new(big.Int).SetBytes(log.Data)
		usdt := new(big.Float).Quo(
			new(big.Float).SetInt(value),
			scale,
		).Text('f', 18)
		isStuck := tx.ContractAddress != w.Address

		return usdt, isStuck
	}
	return "", false
}

func (w *ERC20[T]) IsTransactionMatch(address string, txid string) (string, bool) {
	Tx, err := w.client.TransactionReceipt(txid)
	if err != nil {
		logger.Log.Debugf("[%s]: %s", w.Name, err)
		return "", false
	}

	_Tx, err := w.client.TransactionByHash(txid)
	if err == nil {
		Tx.ContractAddress = *_Tx.To()
	}

	usdt, isStuck := w.Checklogs(Tx, address)
	return usdt, isStuck
}

func (w *ERC20[T]) GetLatestTx(address string) (string, bool) {
	hash, err := w.client.GetERC20(w.Address.Hex(), address)
	if err != nil {
		logger.Log.Warnf("[%s]: %v", w.Name, err)
		return "", false
	}
	return w.IsTransactionMatch(address, hash)
}
