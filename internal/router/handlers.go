package router

import (
	"encoding/json"
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/chat"
	"github.com/ayushgpt01/chatRoomGo/utils"
)

func handleAPIRoutes(mux *http.ServeMux, chatService *chat.ChatService, authService *auth.AuthService) {
	mux.Handle("POST /auth/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload auth.LoginPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		res, err := authService.HandleLogin(r.Context(), payload)
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

	mux.Handle("POST /auth/signup", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload auth.SignupPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		res, err := authService.HandleSignup(r.Context(), payload)
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

	mux.Handle("POST /auth/refresh", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			RefreshToken string `json:"refreshToken"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		token, err := authService.HandleRefresh(r.Context(), payload.RefreshToken)
		// TODO - Add special error types to detect special auth error like token expired etc..
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		var response = struct {
			Token string `json:"token"`
		}{Token: token}

		err = utils.Encode(w, r, http.StatusOK, response)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	}))

	AddProtectedRoutes(mux, chatService, authService)
}
