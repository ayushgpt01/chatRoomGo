package middleware

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/logger"
	"github.com/google/uuid"
)

type contextKey string

const RequestIDKey contextKey = "requestID"

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Add request ID to response headers for client-side tracking
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware provides comprehensive request/response logging
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		// Get request ID from context
		requestID := "unknown"
		if reqID, ok := r.Context().Value(RequestIDKey).(string); ok {
			requestID = reqID
		}

		// Create logger with request context
		log := logger.WithRequestID(requestID)

		// Log incoming request details
		log.Debug("incoming_request",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"content_type", r.Header.Get("Content-Type"),
		)

		// Read and log request body if it exists and logging is enabled
		if r.Body != nil && r.ContentLength > 0 && logger.ShouldLogPayload() {
			body, err := readRequestBody(r)
			if err == nil && len(body) > 0 {
				var jsonBody any

				if json.Unmarshal(body, &jsonBody) == nil {
					log.Debug("request_payload", "payload", jsonBody)
				} else {
					log.Debug("request_payload", "payload", string(body))
				}

			}
		}

		// Wrap response writer to capture status code and response size
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default status code
		}

		// Process the request
		next.ServeHTTP(wrapped, r)

		// Calculate duration
		duration := time.Since(startTime)

		// Log response details
		log.Debug("outgoing_response",
			"status_code", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"duration", duration.String(),
			"response_size", wrapped.size,
		)

		// Log response body if it's small enough and logging is enabled
		if wrapped.size > 0 && wrapped.size < 1024 && logger.ShouldLogPayload() {
			if len(wrapped.body) > 0 {
				var jsonBody any

				if json.Unmarshal(wrapped.body, &jsonBody) == nil {
					log.Debug("response_payload", "payload", jsonBody)
				} else {
					log.Debug("response_payload", "payload", string(wrapped.body))
				}
			}
		}

		// Log request completion summary
		log.Info("request_completed",
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code and response size
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
	body       []byte
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.size += n

	// Capture response body for logging (limit to 1KB to avoid memory issues)
	if rw.size <= 1024 {
		if len(rw.body) == 0 {
			rw.body = make([]byte, 0, 1024)
		}
		rw.body = append(rw.body, b...)
	}

	return n, err
}

// Hijack implements http.Hijacker interface to support WebSocket upgrades
func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	// Check if the underlying ResponseWriter implements Hijacker
	if hijacker, ok := rw.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}

	// If it doesn't implement Hijacker, return an error
	return nil, nil, fmt.Errorf("response writer does not implement http.Hijacker")
}

// readRequestBody reads the request body for logging and restores it
func readRequestBody(r *http.Request) ([]byte, error) {
	// Only read body if it's reasonable size
	if r.ContentLength > 1024*1024 { // 1MB limit
		return nil, nil
	}

	// Read the body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Restore the body for subsequent handlers
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}
