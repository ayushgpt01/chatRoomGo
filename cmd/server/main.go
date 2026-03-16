package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/logger"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func main() {
	logger.Init()

	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, relying on OS environment variables")
	}

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	logger.Info("Starting ChatRoom server", "host", host, "port", port)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		<-ctx.Done()
		logger.Info("Shutting down server...")
	}()

	db, err := sql.Open("pgx", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Error("Failed to open database connection", "error", err)
		os.Exit(1)
	}

	handler := handlerInit(ctx, db)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: handler,
	}

	go func() {
		logger.Info("Server listening", "address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error listening to server", "error", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Error("Error shutting down HTTP server", "error", err)
		}
		if err := db.Close(); err != nil {
			logger.Error("Error closing database connection", "error", err)
		}
	})

	wg.Wait()
}
