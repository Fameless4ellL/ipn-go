package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DBURL                  string
	BotToken               string
	ChatId                 string
	Port                   int
	Verbose                bool
	BalanceTolerance       = 0.01
	ETHapiKey              string
	SOLapiKey              string
	TG_BASE_URL            string
	TELEGRAM_AUTH_BASE_URL string
	TG_WEBHOOK_PATH        string
)

func Init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Read environment variables
	DBURL = getEnv("DB_URL", "")
	BotToken = getEnv("BOT_TOKEN", "")
	ChatId = getEnv("TELEGRAM_CHAT_ID", "")
	TG_BASE_URL = getEnv("TELEGRAM_BASE_URL", "https://api.telegram.org")
	TELEGRAM_AUTH_BASE_URL = getEnv("TELEGRAM_AUTH_BASE_URL", "")
	TG_WEBHOOK_PATH = getEnv("TELEGRAM_WEBHOOK_PATH", "")
	Port, _ = strconv.Atoi(getEnv("PORT", "8080"))
	Verbose, _ = strconv.ParseBool(getEnv("VERBOSE", "false"))
	BalanceTolerance, _ = strconv.ParseFloat(getEnv("BALANCE_TOLERANCE", "0.01"), 64)
	ETHapiKey = getEnv("ETHERSCAN_API_KEY", "")
	SOLapiKey = getEnv("SOLANA_API_KEY", "")
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
