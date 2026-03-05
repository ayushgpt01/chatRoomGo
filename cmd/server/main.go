package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

const HOST = "localhost"
const PORT = "8080"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		<-ctx.Done()
		log.Printf("Shutting down...\n")
	}()

	db, err := sql.Open("sqlite", "my_app_database.db?_fk=1")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error opening db connection: %s\n", err)
		os.Exit(1)
	}

	handler := handlerInit(ctx, db)

	httpServer := &http.Server{
		Addr:    net.JoinHostPort(HOST, PORT),
		Handler: handler,
	}

	go func() {
		log.Printf("Listening on: %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening to server: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Go(func() {
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
		if err := db.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down database connection: %s\n", err)
		}
	})

	wg.Wait()
}
