package auth

import (
	"encoding/json"
	"log"
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
			log.Printf("POST /auth/login - %v\n", err)
			if err.Error() == "invalid credentials" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
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
		// TODO - Add special error types to detect special auth error like token expired etc..
		if err != nil {
			log.Printf("POST /auth/signup - %v\n", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
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
		// TODO - Add special error types to detect special auth error like token expired etc..
		if err != nil {
			log.Printf("POST /auth/refresh - %v\n", err)
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
		// TODO - Add special error types to detect special auth error like token expired etc..
		if err != nil {
			log.Printf("GET /auth/me - %v\n", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
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
		// TODO - Add special error types to detect special auth error like token expired etc..
		if err != nil {
			log.Printf("POST /auth/logout - %v\n", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})
}
