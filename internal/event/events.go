package event

import (
	"errors"
	"log"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

func NewErrorEvent(err error) models.ChatEvent {
	code, message := mapErrorCode(err)
	log.Printf("[Error Event]: %v", err)

	return &models.BaseEvent{
		EventType: models.EventError,
		Data: map[string]any{
			"message": message,
			"code":    code,
		},
	}
}

func mapErrorCode(err error) (code string, message string) {
	switch {
	case errors.Is(err, ErrInvalidPayload):
		return "invalid_payload", err.Error()
	case errors.Is(err, ErrNotRoomMember):
		return "not_room_member", err.Error()
	case errors.Is(err, ErrForbidden):
		return "forbidden", err.Error()
	case errors.Is(err, ErrUnsupportedEvent):
		return "unsupported_event", err.Error()
	default:
		return "internal_error", "something went wrong"
	}
}
