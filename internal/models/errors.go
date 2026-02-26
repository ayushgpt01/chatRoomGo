package models

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUserNotInRoom = errors.New("user is not a member of this room")
	ErrAlreadyInRoom = errors.New("user is already in this room")
)
