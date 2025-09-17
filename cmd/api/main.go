package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go-blocker/internal/config"
	"go-blocker/internal/database"
	logger "go-blocker/internal/log"
	"go-blocker/internal/payment"
	"go-blocker/internal/server"
	"go-blocker/internal/storage"
	blocker "go-blocker/internal/worker/blocker"
	worker "go-blocker/internal/worker/payment"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Listen for the interrupt signal.
	<-ctx.Done()

	logger.Log.Infoln("shutting down gracefully, press Ctrl+C again to force")
	stop()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		logger.Log.Infof("Server forced to shutdown with error: %v", err)
	}

	logger.Log.Infoln("Server exiting")
	done <- true
}

func main() {
	config.Init()
	server := server.NewServer()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine
	go gracefulShutdown(server, done)

	db := database.New()
	repo := database.NewPaymentRepository(db)
	service := payment.NewPaymentService(repo)

	storage.InitStores() // Initialize the global stores, including payment
	blocker.Start(service)
	worker.Start(service)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Log.Panicf("http server error: %s", err)
	}

	<-done
	logger.Log.Infoln("Graceful shutdown complete.")
}
