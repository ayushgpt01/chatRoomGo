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

	"github.com/ayushgpt01/chatRoomGo/internal/chat"
	"github.com/ayushgpt01/chatRoomGo/internal/message"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/router"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
	"github.com/ayushgpt01/chatRoomGo/internal/ws"
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

	userStore, err := user.NewSqliteUserRepo(ctx, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising user repo: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Initialised User Repo\n")

	roomStore, err := room.NewSQLiteRoomRepo(ctx, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising room repo: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Initialised Room Repo\n")

	messageStore, err := message.NewSQLiteMessageRepo(ctx, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising message repo: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Initialised Message Repo\n")

	roomMemberStore, err := chat.NewSQLiteRoomMemberRepo(ctx, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising room member repo: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Initialised Room member Repo\n")

	hub := ws.NewHub(ctx)
	go hub.Cleanup()

	chatService := chat.NewChatService(userStore, roomStore, messageStore, roomMemberStore)
	wsHandler := ws.NewWSHandler(hub, chatService)
	handler := router.HandleRoutes(wsHandler, chatService)

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
