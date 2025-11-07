package main

import (
	"context"
	"net/http"

	application "go-blocker/internal/application/payment"
	repository "go-blocker/internal/infrastructure/payment"
	blockchain "go-blocker/internal/infrastructure/provider"
	storage "go-blocker/internal/infrastructure/storage"
	"go-blocker/internal/infrastructure/telegram"
	worker "go-blocker/internal/infrastructure/worker"
	database "go-blocker/internal/interface/db"
	server "go-blocker/internal/interface/http"
	handler "go-blocker/internal/interface/http/handler"
	"go-blocker/internal/pkg/config"
	logger "go-blocker/internal/pkg/log"
	"go-blocker/internal/rpc"
	"os/signal"
	"syscall"
	"time"
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
	logger.Init()

	// Create a done channel to signal when the shutdown is complete
	done := make(chan bool, 1)

	// Run graceful shutdown in a separate goroutine

	telegram.Init()
	box := storage.NewAddressStore()

	db := database.New()
	repo := repository.NewRepository(db)
	manager := rpc.NewManager()
	watcher := blockchain.NewCurrencyWatcherRegistry(manager)

	service := application.NewService(repo, manager, watcher, box)
	h := handler.NewRepository(service)

	router := server.RegisterRoutes(h)
	srv := server.NewServer(router)
	go gracefulShutdown(srv, done)

	work := worker.NewWorker(service, 10*time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go work.Start(ctx)

	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logger.Log.Panicf("http server error: %s", err)
	}

	<-done
	logger.Log.Infoln("Graceful shutdown complete.")
}
