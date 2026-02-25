package chat

import (
	"log"

	"github.com/ayushgpt01/chatRoomGo/internal/types"
)

func NewErrorEvent(err error) types.ChatEvent {
	code, message := mapErrorCode(err)
	log.Printf("[Error Event]: %s", err)

	return &types.BaseEvent{
		EventType: types.EventError,
		Data: map[string]any{
			"message": message,
			"code":    code,
		},
	}
}

func mapErrorCode(err error) (code string, message string) {
	switch err {
	case ErrInvalidPayload:
		return "invalid_payload", err.Error()
	case ErrNotRoomMember:
		return "not_room_member", err.Error()
	case ErrForbidden:
		return "forbidden", err.Error()
	case ErrUnsupportedEvent:
		return "unsupported_event", err.Error()
	default:
		return "internal_error", "something went wrong"
	}
}
