package main

import (
	"context"
	"database/sql"
	"net/http"
	"os"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/event"
	"github.com/ayushgpt01/chatRoomGo/internal/logger"
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
		logger.Error("Failed to initialize user repo", "error", err)
		os.Exit(1)
	}

	logger.Info("Initialized User Repo")

	roomStore, err := room.NewSQLiteRoomRepo(ctx, db)
	if err != nil {
		logger.Error("Failed to initialize room repo", "error", err)
		os.Exit(1)
	}
	logger.Info("Initialized Room Repo")

	messageStore, err := message.NewSQLiteMessageRepo(ctx, db)
	if err != nil {
		logger.Error("Failed to initialize message repo", "error", err)
		os.Exit(1)
	}
	logger.Info("Initialized Message Repo")

	roomMemberStore, err := room.NewSQLiteRoomMemberRepo(ctx, db)
	if err != nil {
		logger.Error("Failed to initialize room member repo", "error", err)
		os.Exit(1)
	}
	logger.Info("Initialized Room member Repo")

	authStore, err := auth.NewSQLiteAuthRepo(ctx, db)
	if err != nil {
		logger.Error("Failed to initialize auth repo", "error", err)
		os.Exit(1)
	}
	logger.Info("Initialized Auth Repo")

	if err := seed.SeedChatData(context.Background(), db); err != nil {
		logger.Error("Failed to seed chat data", "error", err)
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
