package router

import (
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/chat"
)

func handleAPIRoutes(mux *http.ServeMux, chatService *chat.ChatService) {
	mux.Handle("POST /users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	mux.Handle("POST /rooms", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	mux.Handle("POST /rooms/join", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	mux.Handle("GET /rooms/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
}
