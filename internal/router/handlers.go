package router

import (
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/chat"
)

func handleAPIRoutes(mux *http.ServeMux, chatService *chat.ChatService, authService *auth.AuthService) {
	apiMux := http.NewServeMux()

	// Public Routes
	apiMux.Handle("POST /auth/login", auth.HandleLogin(authService))
	apiMux.Handle("POST /auth/signup", auth.HandleSignup(authService))
	apiMux.Handle("POST /auth/refresh", auth.HandleRefresh(authService))

	protectedMux := http.NewServeMux()

	// Protected Routes
	protectedMux.Handle("GET /auth/me", auth.HandleMe(authService))
	protectedMux.Handle("POST /auth/logout", auth.HandleLogout(authService))
	// protectedMux.Handle("POST /rooms", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	// protectedMux.Handle("POST /rooms/join", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	// protectedMux.Handle("GET /rooms/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	apiMux.Handle("/", authService.Middleware(protectedMux))

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))
}
