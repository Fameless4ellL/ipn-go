package watcher

import (
	"context"
	"fmt"
	app "go-blocker/internal/application/payment"
	"go-blocker/internal/deprecated/storage"
	domain "go-blocker/internal/domain/payment"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/rpc"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
)

type USDT struct {
	S *app.Service
}

func (w *USDT) Name() string {
	return "USDT"
}

func (w *USDT) Chain() string {
	return "ethereum"
}

func (w *USDT) Address() common.Address {
	return common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7")
}

func (w *USDT) Decimals() int {
	return int(big.NewInt(1_000_000).Int64())
}

func (w *USDT) Abi() abi.ABI {
	abiJSON := `[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_upgradedAddress","type":"address"}],"name":"deprecate","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"deprecated","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_evilUser","type":"address"}],"name":"addBlackList","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"upgradedAddress","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balances","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"maximumFee","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"_totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"unpause","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_maker","type":"address"}],"name":"getBlackListStatus","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"allowed","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"paused","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"who","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"pause","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"getOwner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newBasisPoints","type":"uint256"},{"name":"newMaxFee","type":"uint256"}],"name":"setParams","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"uint256"}],"name":"issue","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"uint256"}],"name":"redeem","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"_owner","type":"address"},{"name":"_spender","type":"address"}],"name":"allowance","outputs":[{"name":"remaining","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"basisPointsRate","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"isBlackListed","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_clearedUser","type":"address"}],"name":"removeBlackList","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"MAX_UINT","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_blackListedUser","type":"address"}],"name":"destroyBlackFunds","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"_initialSupply","type":"uint256"},{"name":"_name","type":"string"},{"name":"_symbol","type":"string"},{"name":"_decimals","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"}],"name":"Issue","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"amount","type":"uint256"}],"name":"Redeem","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"newAddress","type":"address"}],"name":"Deprecate","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"feeBasisPoints","type":"uint256"},{"indexed":false,"name":"maxFee","type":"uint256"}],"name":"Params","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_blackListedUser","type":"address"},{"indexed":false,"name":"_balance","type":"uint256"}],"name":"DestroyedBlackFunds","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_user","type":"address"}],"name":"AddedBlackList","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_user","type":"address"}],"name":"RemovedBlackList","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"},{"indexed":true,"name":"spender","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Approval","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[],"name":"Pause","type":"event"},{"anonymous":false,"inputs":[],"name":"Unpause","type":"event"}]`
	parsedABI, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		log.Fatal(err)
	}
	return parsedABI
}

func (w *USDT) HasActiveAddresses() bool {
	return storage.PaymentAddressStore.Len() > 0
}

func (w *USDT) CheckTransactions(m *rpc.Manager, client *ethclient.Client, block []*types.Receipt) (uuid.UUID, error) {
	for _, tx := range block {
		w.Checklogs(client, tx)
	}

	return uuid.Nil, fmt.Errorf("no matching address found")

}

func (w *USDT) CheckAddress(address string) (uuid.UUID, error) {
	// Base Transfer From - To
	if id, ok := storage.PaymentAddressStore.Get(address); ok {
		return id, nil
	}

	return uuid.Nil, fmt.Errorf("no matching address found")
}

func (w *USDT) Checklogs(client *ethclient.Client, tx *types.Receipt) {
	for _, log := range tx.Logs {
		if len(log.Topics) >= 3 {
			toTopic := log.Topics[2]
			if len(toTopic.Bytes()) == 32 {
				to := common.BytesToAddress(toTopic.Bytes()[12:])
				id, err := w.CheckAddress(to.String())
				if err != nil {
					continue
				}
				value := new(big.Int).SetBytes(log.Data)
				usdt := new(big.Float).Quo(
					new(big.Float).SetInt(value),
					new(big.Float).SetInt(big.NewInt(1_000_000)),
				).Text('f', 18)
				txId := log.TxHash.Hex()
				isContractMatch := tx.ContractAddress != w.Address()

				w.S.Status(id, domain.Completed, &usdt, &txId, &isContractMatch)
			}
		}
	}
}

func (w *USDT) GetPendingBalance(client *ethclient.Client, wallet common.Address) big.Float {
	abi := w.Abi()

	data, err := abi.Pack("balanceOf", wallet)
	if err != nil {
		log.Fatal(err)
	}

	ERC20Address := w.Address()
	msg := ethereum.CallMsg{
		To:   &ERC20Address,
		Data: data,
	}

	result, err := client.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Fatal(err)
	}

	var balance *big.Int
	err = abi.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return *big.NewFloat(0)
	}
	usdtBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(float64(w.Decimals())))

	if usdtBalance.Cmp(big.NewFloat(0)) > 0 {
		logger.Log.Infof("USDT: Pending balance for address %s: %f", wallet.Hex(), usdtBalance)
		id, err := w.CheckAddress(wallet.Hex())
		if err != nil {
			return *usdtBalance
		}
		w.S.Status(id, domain.Received, nil, nil, nil)
		storage.PaymentAddressStore.SetPending(wallet.Hex(), true)
	}

	return *usdtBalance
}

func (w *USDT) IsTransactionMatch(client *ethclient.Client, tx *app.CheckTxRequest) (*types.Transaction, string, bool) {
	Tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(tx.TxID))
	if err != nil {
		logger.Log.Debugf("ETH: %s", err)
	}

	wei := Tx.Value()
	eth := new(big.Float).Quo(new(big.Float).SetInt(wei), big.NewFloat(1e18)).Text('f', 18)
	logger.Log.Infof("ETH: Incoming transaction %s, Amount: %s", Tx.Hash().Hex(), eth)

	return Tx, eth, true
}
