# Go-Blocker

Go-Blocker is a backend service for monitoring blockchain payments, supporting ETH, USDT, and USDC on Ethereum. It provides REST API endpoints to check and find transactions, and sends notifications via Telegram.

## Features

- REST API for payment status and transaction lookup
- Supports ETH, USDT, USDC on Ethereum
- Swagger documentation
- Docker support

## Getting Started

### Prerequisites

- Go 1.20+
- Docker (optional)
- Telegram bot (for notifications)
- Etherscan API key

### Installation

1. Set up your `.env` file:
   ```env
   PORT=8080
   GIN_MODE=release
   APP_ENV=local
   VERBOSE=true
   DB_URL=test.db
   BOT_TOKEN=your_telegram_bot_token
   TELEGRAM_CHAT_ID=your_chat_id
   BALANCE_TOLERANCE=0.01
   ETHERSCAN_API_KEY=your_etherscan_api_key
   ```

2. Run the server:
   ```sh
   go run cmd/api/main.go
   ```

## Environment Variables

| Variable             | Description                                                                                   | Example Value                          |
|----------------------|-----------------------------------------------------------------------------------------------|----------------------------------------|
| `PORT`               | The port on which the HTTP server will listen.                                                | `8080`                                 |
| `GIN_MODE`           | Gin web framework mode (`release` for production, `debug` for development).                   | `release`                              |
| `APP_ENV`            | Application environment (e.g., `local`, `production`).                                       | `local`                                |
| `VERBOSE`            | Enable verbose logging (`true` for debug logs, `false` for info only).                       | `true`                                 |
| `DB_URL`             | Path or connection string for the SQLite database file.                                       | `test.db`                              |
| `BOT_TOKEN`          | Telegram bot token for sending payment notifications.                                         | `8223641937:AAEkz3vyLDn9jJ8Lh8THiC5-MFiaanhjUEg` |
| `TELEGRAM_CHAT_ID`   | Telegram chat ID where notifications will be sent.                                            | `1792255940`                           |
| `BALANCE_TOLERANCE`  | Allowed tolerance for payment amount mismatch (float, in token units).                        | `0.01`                                 |
| `ETHERSCAN_API_KEY`  | API key for accessing Etherscan to fetch transaction data.                                    | `VJ68G29Y3EJBB4QYGEEEX1DK9US27N6AAB`   |

### docs
```sh
swag init -g cmd/api/main.go
```