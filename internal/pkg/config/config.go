package config

import (
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
	ETHapiKey        string
	SOLapiKey        string
)

func Init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
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
	// SOLapiKey := os.Getenv("SOLANA_API_KEY")
}
