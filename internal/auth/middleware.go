package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/ayushgpt01/chatRoomGo/internal/logger"
	"github.com/ayushgpt01/chatRoomGo/internal/middleware"
)

type contextKey string

const UserIDKey contextKey = "userId"

func (srv *AuthService) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get request ID for context
		requestID := "unknown"
		if reqID, ok := r.Context().Value(middleware.RequestIDKey).(string); ok {
			requestID = reqID
		}
		log := logger.WithRequestID(requestID)

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Warn("auth_failed_missing_token",
				"reason", "missing_or_invalid_token",
				"auth_header_present", authHeader != "",
			)
			http.Error(w, "Unauthorized: Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Use your existing service method!
		userId, err := srv.getByAccessToken(tokenString)
		if err != nil {
			log.Warn("auth_failed_invalid_token",
				"reason", "invalid_token",
				"error", err.Error(),
				"token_length", len(tokenString),
			)
			// Return 401 so the frontend Axios interceptor triggers a refresh
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		log.Info("auth_success",
			"user_id", userId,
			"method", r.Method,
			"path", r.URL.Path,
		)

		ctx := context.WithValue(r.Context(), UserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (srv *AuthService) OptionalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// If no header, just move to the next handler with empty context
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userId, err := srv.getByAccessToken(tokenString)
		if err != nil {
			// Return 401 so the frontend Axios interceptor triggers a refresh
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
