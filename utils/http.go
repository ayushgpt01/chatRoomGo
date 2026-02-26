package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
}

// Validator is an object that can be validated.
type Validator interface {
	// Valid checks the object and returns any
	// problems. If len(problems) == 0 then
	// the object is valid.
	Valid(ctx context.Context) (problems map[string]string)
}

func DecodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}
	if problems := v.Valid(r.Context()); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}
	return v, nil, nil
}

func HandleDecode[T Validator](w http.ResponseWriter, r *http.Request) (T, bool) {
	temp, problems, err := DecodeValid[T](r)

	if err != nil {
		if len(problems) > 0 {
			Encode(w, r, http.StatusUnprocessableEntity, map[string]any{"errors": problems})
			return temp, false
		}
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return temp, false
	}

	return temp, true
}

func HandleServiceError(w http.ResponseWriter, path string, err error) {
	log.Printf("%s - %v\n", path, err)

	// Switch on error types or messages
	switch {
	case errors.Is(err, models.ErrUnauthorized):
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

	case errors.Is(err, models.ErrNotFound):
		http.Error(w, "Resource not found", http.StatusNotFound)

	case errors.Is(err, models.ErrUserNotInRoom):
		http.Error(w, err.Error(), http.StatusForbidden)

	case errors.Is(err, models.ErrAlreadyInRoom):
		http.Error(w, err.Error(), http.StatusConflict)

	default:
		http.Error(w, "Internal Server error", http.StatusInternalServerError)
	}
}
