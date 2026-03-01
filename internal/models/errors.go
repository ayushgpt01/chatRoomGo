package models

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal error")
)
