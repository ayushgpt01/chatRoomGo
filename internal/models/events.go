package models

import (
	"encoding/json"
)

type OutgoingEvent struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

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

// Incoming events are named like commands
type IncomingEventType string
type OutgoingEventType string

const (
	EventSendMessage   IncomingEventType = "message.send"
	EventJoinRoom      IncomingEventType = "room.join"
	EventLeaveRoom     IncomingEventType = "room.leave"
	EventEditMessage   IncomingEventType = "message.edit"
	EventDeleteMessage IncomingEventType = "message.delete"
)

const (
	EventMessageCreated OutgoingEventType = "message_created"
	EventMessageUpdated OutgoingEventType = "message_updated"
	EventMessageDeleted OutgoingEventType = "message_deleted"
	EventUserJoinedRoom OutgoingEventType = "user_joined_room"
	EventUserLeftRoom   OutgoingEventType = "user_left_room"

	EventError OutgoingEventType = "error"
)

// EventMessageCreated - "message_created"
type MessageCreatedPayload struct {
	Message *ResponseMessage `json:"message"`
}

type MessageCreatedEvent struct {
	Data MessageCreatedPayload
}

func (e *MessageCreatedEvent) Type() string {
	return string(EventMessageCreated)
}

func (e *MessageCreatedEvent) Payload() any {
	return e.Data
}

// EventMessageUpdated - "message_updated"
type MessageUpdatedPayload struct {
	Message *ResponseMessage `json:"message"`
}

type MessageUpdatedEvent struct {
	Data MessageUpdatedPayload
}

func (e *MessageUpdatedEvent) Type() string {
	return string(EventMessageUpdated)
}

func (e *MessageUpdatedEvent) Payload() any {
	return e.Data
}

// EventMessageDeleted - "message_deleted"
type MessageDeletedPayload struct {
	MessageID MessageId `json:"messageId"`
	RoomID    RoomId    `json:"roomId"`
}

type MessageDeletedEvent struct {
	Data MessageDeletedPayload
}

func (e *MessageDeletedEvent) Type() string {
	return string(EventMessageDeleted)
}

func (e *MessageDeletedEvent) Payload() any {
	return e.Data
}

// EventUserJoinedRoom - "user_joined_room"
type UserJoinedRoomPayload struct {
	RoomID RoomId `json:"roomId"`
	UserID UserId `json:"userId"`
}

type UserJoinedRoomEvent struct {
	Data UserJoinedRoomPayload
}

func (e *UserJoinedRoomEvent) Type() string {
	return string(EventUserJoinedRoom)
}

func (e *UserJoinedRoomEvent) Payload() any {
	return e.Data
}

// EventUserLeftRoom - "user_left_room"
type UserLeftRoomPayload struct {
	RoomID RoomId `json:"roomId"`
	UserID UserId `json:"userId"`
}

type UserLeftRoomEvent struct {
	Data UserLeftRoomPayload
}

func (e *UserLeftRoomEvent) Type() string {
	return string(EventUserLeftRoom)
}

func (e *UserLeftRoomEvent) Payload() any {
	return e.Data
}

// EventError - "error"
type ErrorPayload struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ErrorEvent struct {
	Data ErrorPayload
}

func (e *ErrorEvent) Type() string {
	return string(EventError)
}

func (e *ErrorEvent) Payload() any {
	return e.Data
}
