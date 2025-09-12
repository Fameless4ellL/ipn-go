package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"

	"go-blocker/internal/config"
	"go-blocker/internal/database"
)

type Server struct {
	port int
	db   *gorm.DB
}

func NewServer() *http.Server {
	NewServer := &Server{
		port: config.Port,
		db:   database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
