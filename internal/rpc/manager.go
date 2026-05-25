package rpc

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"time"

	domain "go-blocker/internal/domain/blockchain"
	logger "go-blocker/internal/pkg/log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"resty.dev/v3"
)

type Domain map[domain.ChainType][]string

var Domains = Domain{
	domain.Ethereum: {
		"https://eth.drpc.org",
		"https://ethereum-rpc.publicnode.com",
	},
	domain.Solana: {
		"https://api.mainnet-beta.solana.com",
		"https://solana.drpc.org",
	},
	domain.Binance: {
		"https://bsc.drpc.org",
	},
}

type Manager struct {
	nodes Domain
}

func NewManager() *Manager {
	return &Manager{nodes: Domains}
}

func (m *Manager) GetClientForChain(chain domain.ChainType) (*domain.Node, string, error) {
	node, _ := m.nodes[chain]

	rr, err := resty.NewRoundRobin(node...)
	if err != nil {
		logger.Log.Warn("Can't create round robin", slog.String("error", err.Error()))
		return nil, node[0], err
	}

	client := resty.New().SetLoadBalancer(rr).SetTimeout(5 * time.Second)
	HttpClient := client.Client()

	rpcClient, err := rpc.DialOptions(context.Background(), node[0], rpc.WithHTTPClient(HttpClient))
	if err != nil {
		logger.Log.Warn("Can't create RPC client", slog.String("error", err.Error()))
		return nil, node[0], err
	}

	return NewNode(rr, ethclient.NewClient(rpcClient)), node[0], nil
}

func NewNode(rr *resty.RoundRobin, client *ethclient.Client) *domain.Node {
	return &domain.Node{Rr: rr, Client: client}
}

func Do[T any](node domain.Node, methodName string, args ...any) (T, error) {
	var zero T

	for i := range 3 {
		v := reflect.ValueOf(node.Client)
		method := v.MethodByName(methodName)

		if !method.IsValid() {
			continue
		}

		inputs := make([]reflect.Value, len(args))
		for i, arg := range args {
			if arg == nil {
				inputs[i] = reflect.Zero(method.Type().In(i))
			} else {
				inputs[i] = reflect.ValueOf(arg)
			}
		}

		results := method.Call(inputs)

		if len(results) == 0 {
			continue
		}

		var decrErr error
		if len(results) > 1 {
			lastVal := results[len(results)-1]

			if !lastVal.IsNil() && lastVal.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
				decrErr = lastVal.Interface().(error)
				logger.Log.Warn("method %s returned an error: %v", methodName, decrErr, slog.Int("attempt", i))

				url, err := node.Rr.Next()
				if err != nil {
					logger.Log.Warn("Can't get next URL from round robin", slog.String("error", err.Error()))
					continue
				}

				client, err := ethclient.Dial(url)
				if err != nil {
					logger.Log.Warn("Can't dial URL", slog.String("url", url), slog.String("error", err.Error()))
					continue
				}
				node.Client = client

				continue
			}
		}

		val, ok := results[0].Interface().(T)
		if !ok {
			if results[0].IsNil() {
				return zero, decrErr
			}
			return zero, fmt.Errorf("method returned %T, but expected %T", results[0].Interface(), zero)
		}

		return val, decrErr

	}

	return zero, fmt.Errorf("method %s failed after 3 attempts", methodName)
}
