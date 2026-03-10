package logger

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var (
	defaultLogger *slog.Logger
)

// Init initializes the default logger with the specified configuration
func Init() {
	level := getLogLevel()

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if os.Getenv("LOG_FORMAT") == "text" {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// getLogLevel determines the log level from environment variable
func getLogLevel() slog.Level {
	envLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch envLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		// Default to INFO for production, DEBUG for development
		if os.Getenv("GO_ENV") == "production" {
			return slog.LevelInfo
		}
		return slog.LevelDebug
	}
}

// SetLogger allows swapping out the default logger
func SetLogger(logger *slog.Logger) {
	defaultLogger = logger
	slog.SetDefault(defaultLogger)
}

// GetLogger returns the current default logger
func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		Init()
	}
	return defaultLogger
}

// WithRequestID adds request ID to logger context
func WithRequestID(requestID string) *slog.Logger {
	return GetLogger().With("request_id", requestID)
}

// WithUserID adds user ID to logger context
func WithUserID(userID int) *slog.Logger {
	return GetLogger().With("user_id", userID)
}

// WithRoomID adds room ID to logger context
func WithRoomID(roomID int) *slog.Logger {
	return GetLogger().With("room_id", roomID)
}

// Info logs an info message with structured context
func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

// Warn logs a warning message with structured context
func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

// Error logs an error message with structured context
func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

// ShouldLogPayload returns true if payload should be logged (based on config)
func ShouldLogPayload() bool {
	return GetLogger().Enabled(context.Background(), slog.LevelDebug)
}

// ParseInt safely parses string to int for logging
func ParseInt(s string) int {
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	return 0
}
