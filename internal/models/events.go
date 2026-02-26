package models

import (
	"encoding/json"
)

type IncomingEvent struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type ChatEvent interface {
	Type() string
	Payload() any
}

type HubBroadcaster interface {
	Broadcast(roomId int64, event ChatEvent) error
}

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

type BaseEvent struct {
	EventType OutgoingEventType
	Data      any
}

func (e *BaseEvent) Type() string {
	return string(e.EventType)
}

func (e *BaseEvent) Payload() any {
	return e.Data
}
