package router

import (
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/logger"
	"github.com/ayushgpt01/chatRoomGo/internal/message"
	"github.com/ayushgpt01/chatRoomGo/internal/middleware"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/ws"
	"github.com/rs/cors"
)

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				logger.Error("Panic recovered in HTTP handler",
					"panic", rec,
					"method", r.Method,
					"path", r.URL.Path,
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func HandleRoutes(wsHandler *ws.Wshandler, authService *auth.AuthService, roomService *room.RoomService, messageService *message.MessageService) http.Handler {
	logger.Info("Setting up routes...")

	mux := http.NewServeMux()

	handleAPIRoutes(mux, authService, roomService, messageService)
	handleViews(mux)
	mux.Handle("/ws", wsHandler)

	var handler http.Handler = mux
	handler = middleware.RequestIDMiddleware(handler)
	handler = recoverMiddleware(handler)
	handler = middleware.LoggingMiddleware(handler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            false,
	})

	return c.Handler(handler)
}
