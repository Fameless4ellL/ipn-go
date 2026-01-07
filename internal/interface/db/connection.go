package database

import (
	"go-blocker/internal/infrastructure/payment"
	"go-blocker/internal/pkg/config"
	logger "go-blocker/internal/pkg/log"
	"log/slog"

	"github.com/glebarez/sqlite" // Pure Go SQLite driver
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

var db *gorm.DB

func New() *gorm.DB {
	if db != nil {
		return db
	}
	// Initialize the database connection
	if config.DBURL == "" {
		logger.Log.Error("DB_URL environment variable is not set")
	}
	// Use SQLite for simplicity, but you can change this to any other database driver
	db, err := gorm.Open(sqlite.Open(config.DBURL), &gorm.Config{})
	if err != nil {
		logger.Log.Error("failed to connect database", slog.String("error", err.Error()))
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&payment.PaymentModel{})
	if err != nil {
		logger.Log.Error("failed to migrate schema", slog.String("error", err.Error()))
	}

	return db
}
