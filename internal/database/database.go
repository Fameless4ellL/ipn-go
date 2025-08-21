package database

import (
	"go-blocker/internal/payment"
	"log"
	"os"

	"github.com/glebarez/sqlite" // Pure Go SQLite driver
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

var (
	dburl = os.Getenv("DB_URL")
	db    *gorm.DB
)

func New() *gorm.DB {
	if db != nil {
		return db
	}
	// Initialize the database connection
	if dburl == "" {
		log.Fatal("DB_URL environment variable is not set")
	}
	// Use SQLite for simplicity, but you can change this to any other database driver
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(&payment.Payment{})
	if err != nil {
		log.Fatal("failed to migrate schema:", err)
	}

	return db
}
