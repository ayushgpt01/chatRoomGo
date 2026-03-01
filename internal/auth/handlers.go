package auth

import (
	"encoding/json"
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/utils"
)

func HandleLogin(srv *AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload LoginPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		res, err := srv.HandleLogin(r.Context(), payload)
		if err != nil {
			utils.HandleServiceError(w, "POST /auth/login", err)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, res)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}

func HandleSignup(srv *AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload SignupPayload

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		res, err := srv.HandleSignup(r.Context(), payload)
		if err != nil {
			utils.HandleServiceError(w, "POST /auth/signup", err)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, res)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}

func HandleRefresh(srv *AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			RefreshToken string `json:"refreshToken"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		token, err := srv.HandleRefresh(r.Context(), payload.RefreshToken)
		if err != nil {
			utils.HandleServiceError(w, "POST /auth/refresh", err)
			return
		}

		var response = struct {
			Token string `json:"token"`
		}{Token: token}

		err = utils.Encode(w, r, http.StatusOK, response)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}

func HandleMe(srv *AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Missing or invalid authorization header", http.StatusUnauthorized)
			return
		}
		accessToken := authHeader[7:]

		res, err := srv.GetCurrentUser(r.Context(), accessToken)
		if err != nil {
			utils.HandleServiceError(w, "POST /auth/me", err)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, res)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}

func HandleLogout(srv *AuthService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			RefreshToken string `json:"refreshToken"`
		}

		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		err := srv.HandleLogout(r.Context(), payload.RefreshToken)
		if err != nil {
			utils.HandleServiceError(w, "POST /auth/logout", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
