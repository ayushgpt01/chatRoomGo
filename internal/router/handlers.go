package router

import (
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
)

func handleAPIRoutes(mux *http.ServeMux, authService *auth.AuthService, roomService *room.RoomService) {
	apiMux := http.NewServeMux()

	// Public Routes
	apiMux.Handle("POST /auth/login", auth.HandleLogin(authService))
	apiMux.Handle("POST /auth/signup", auth.HandleSignup(authService))
	apiMux.Handle("POST /auth/refresh", auth.HandleRefresh(authService))

	// Allows guest as well
	apiMux.Handle("POST /room/join", authService.OptionalMiddleware(room.HandleJoinRoom(roomService)))

	protectedMux := http.NewServeMux()

	// Protected Routes
	protectedMux.Handle("GET /auth/me", auth.HandleMe(authService))
	protectedMux.Handle("POST /auth/logout", auth.HandleLogout(authService))
	protectedMux.Handle("POST /room/leave", room.HandleLeaveRoom(roomService))
	protectedMux.Handle("GET /room/getAll", room.HandleGetRooms(roomService))
	protectedMux.Handle("POST /room/create", room.HandleCreateRoom(roomService))

	apiMux.Handle("/", authService.Middleware(protectedMux))

	mux.Handle("/api/", http.StripPrefix("/api", apiMux))
}
