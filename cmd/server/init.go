package main

import (
	"context"
	"database/sql"
	"net/http"
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
	userStore := user.NewPostgresUserRepo(ctx, db)
	roomStore := room.NewPostgresRoomRepo(ctx, db)
	messageStore := message.NewPostgresMessageRepo(ctx, db)
	roomMemberStore := room.NewPostgresRoomMemberRepo(ctx, db)
	authStore := auth.NewPostgresAuthRepo(ctx, db)

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
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				cleanupCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				authService.HandleCleanup(cleanupCtx)
				cancel()
			}
		}
	}()

	return router.HandleRoutes(wsHandler, authService, roomService, messageService)
}
