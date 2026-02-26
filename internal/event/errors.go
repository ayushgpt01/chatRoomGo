package event

import "errors"

var (
	ErrInvalidPayload   = errors.New("invalid payload")
	ErrNotRoomMember    = errors.New("user is not a member of the room")
	ErrForbidden        = errors.New("forbidden")
	ErrUnsupportedEvent = errors.New("unsupported event type")
)
