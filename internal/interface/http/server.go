package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"

	database "go-blocker/internal/interface/db"
	"go-blocker/internal/pkg/config"
)

type Server struct {
	port int
	db   *gorm.DB
}

func NewServer(h *gin.Engine) *http.Server {
	NewServer := &Server{
		port: config.Port,
		db:   database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      h,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
