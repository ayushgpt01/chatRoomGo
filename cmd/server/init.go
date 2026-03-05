package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/event"
	"github.com/ayushgpt01/chatRoomGo/internal/message"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/router"
	"github.com/ayushgpt01/chatRoomGo/internal/seed"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
	"github.com/ayushgpt01/chatRoomGo/internal/ws"
)

func handlerInit(ctx context.Context, db *sql.DB) http.Handler {
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

	roomMemberStore, err := room.NewSQLiteRoomMemberRepo(ctx, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising room member repo: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Initialised Room member Repo\n")

	authStore, err := auth.NewSQLiteAuthRepo(ctx, db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error initialising auth repo: %s\n", err)
		os.Exit(1)
	}
	log.Printf("Initialised Auth Repo\n")

	if err := seed.SeedChatData(context.Background(), db); err != nil {
		log.Fatal(err)
	}

	hub := ws.NewHub(ctx)
	go hub.Cleanup()

	authService := auth.NewAuthService(userStore, authStore)
	eventService := event.NewEventService(userStore, roomStore, messageStore, roomMemberStore)
	roomService := room.NewRoomService(roomMemberStore, roomStore, authService, hub)
	messageService := message.NewMessageService(messageStore, roomMemberStore, hub)
	wsHandler := ws.NewWSHandler(hub, eventService)

	// Cleanup expired tokens every 1 hour
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			authService.HandleCleanup(context.Background())
		}
	}()

	return router.HandleRoutes(wsHandler, authService, roomService, messageService)
}
