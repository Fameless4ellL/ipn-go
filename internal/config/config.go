package config

import (
	"go-blocker/internal/rpc"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DBURL            string
	BotToken         string
	ChatId           string
	Port             int
	APIKey           string
	Verbose          bool
	BalanceTolerance = 0.01
	Nodes            []rpc.RPCNode
	ETHapiKey        string
)

func Init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	Nodes = []rpc.RPCNode{
		// {URL: "https://eth.drpc.org", Chain: rpc.Ethereum, Healthy: true},                       // has trace_block
		{URL: "https://api.noderpc.xyz/rpc-mainnet/public", Chain: rpc.Ethereum, Healthy: true}, // has trace_block
		{URL: "https://ethereum-public.nodies.app", Chain: rpc.Ethereum, Healthy: true},         // has trace_block
		// {URL: "https://endpoints.omniatech.io/v1/eth/mainnet/public", Chain: rpc.Ethereum, Healthy: true, Processing: false}, // has trace_block
		// {URL: "https://eth.api.onfinality.io/public", Chain: rpc.Ethereum, Healthy: true, Processing: false}, // has trace_block
		// {URL: "https://eth.llamarpc.com", Chain: rpc.Ethereum, Healthy: true}, not trace_block
		// {URL: "https://ethereum-rpc.publicnode.com", Chain: rpc.Ethereum, Healthy: true}, // not trace_block
		// {URL: "https://go.getblock.io/aefd01aa907c4805ba3c00a9e5b48c6b", Chain: rpc.Ethereum, Healthy: true}, too many requests and no support for trace_block
		// {URL: "https://sepolia.drpc.org", Chain: rpc.Ethereum, Healthy: true}, // test
	}

	// Read environment variables
	DBURL = os.Getenv("DB_URL")
	APIKey = os.Getenv("API_KEY")
	BotToken = os.Getenv("BOT_TOKEN")
	ChatId = os.Getenv("TELEGRAM_CHAT_ID")
	Port, _ = strconv.Atoi(os.Getenv("PORT"))
	Verbose, _ = strconv.ParseBool(os.Getenv("VERBOSE"))
	BalanceTolerance, _ = strconv.ParseFloat(os.Getenv("BALANCE_TOLERANCE"), 64)
	ETHapiKey = os.Getenv("ETHERSCAN_API_KEY")
}
