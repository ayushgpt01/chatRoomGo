package router

import (
	"encoding/json"
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/chat"
	"github.com/ayushgpt01/chatRoomGo/utils"
)

func AddProtectedRoutes(mux *http.ServeMux, chatService *chat.ChatService, authService *auth.AuthService) {
	protectedMux := http.NewServeMux()

	protectedMux.Handle("GET /auth/me", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Missing or invalid authorization header", http.StatusUnauthorized)
			return
		}
		accessToken := authHeader[7:]

		res, err := authService.GetCurrentUser(r.Context(), accessToken)
		// TODO - Add special error types to detect special auth error like token expired etc..
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, res)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	}))

	protectedMux.Handle("POST /auth/logout", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			RefreshToken string `json:"refreshToken"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		err := authService.HandleLogout(r.Context(), payload.RefreshToken)
		// TODO - Add special error types to detect special auth error like token expired etc..
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}))

	// protectedMux.Handle("POST /rooms", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	// protectedMux.Handle("POST /rooms/join", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	// protectedMux.Handle("GET /rooms/{id}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	mux.Handle("/", authService.Middleware(protectedMux))
}
