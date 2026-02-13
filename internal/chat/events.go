package chat

import "log"

type IncomingEventType string

const (
	EventSendMessage   IncomingEventType = "send_message"
	EventJoinRoom      IncomingEventType = "join_room"
	EventLeaveRoom     IncomingEventType = "leave_room"
	EventEditMessage   IncomingEventType = "edit_message"
	EventDeleteMessage IncomingEventType = "delete_message"
)

type OutgoingEventType string

const (
	EventMessageCreated OutgoingEventType = "message_created"
	EventMessageUpdated OutgoingEventType = "message_updated"
	EventMessageDeleted OutgoingEventType = "message_deleted"
	EventUserJoinedRoom OutgoingEventType = "user_joined_room"
	EventUserLeftRoom   OutgoingEventType = "user_left_room"

	EventError OutgoingEventType = "error"
)

type ChatEvent interface {
	Type() string
	Payload() any
}

type BaseEvent struct {
	eventType OutgoingEventType
	payload   any
}

func (e *BaseEvent) Type() string {
	return string(e.eventType)
}

func (e *BaseEvent) Payload() any {
	return e.payload
}

func NewErrorEvent(err error) ChatEvent {
	code, message := mapErrorCode(err)
	log.Printf("[Error Event]: %s", err)

	return &BaseEvent{
		eventType: EventError,
		payload: map[string]any{
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
