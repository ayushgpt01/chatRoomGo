package event

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ayushgpt01/chatRoomGo/internal/message"
	"github.com/ayushgpt01/chatRoomGo/internal/models"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
)

type eventHandler func(
	ctx context.Context,
	roomID models.RoomId,
	userID models.UserId,
	data models.IncomingEvent,
) (models.ChatEvent, error)

type EventService struct {
	userStore       user.UserStore
	roomStore       room.RoomStore
	messageStore    message.MessageStore
	roomMemberStore room.RoomMemberStore

	handlers map[models.IncomingEventType]eventHandler
}

func NewEventService(
	userStore user.UserStore,
	roomStore room.RoomStore,
	messageStore message.MessageStore,
	roomMemberStore room.RoomMemberStore,
) *EventService {
	srv := &EventService{
		userStore:       userStore,
		roomStore:       roomStore,
		messageStore:    messageStore,
		roomMemberStore: roomMemberStore,
		handlers:        make(map[models.IncomingEventType]eventHandler),
	}

	srv.handlers[models.EventJoinRoom] = srv.handleJoinRoom
	srv.handlers[models.EventLeaveRoom] = srv.handleLeaveRoom
	srv.handlers[models.EventSendMessage] = srv.handleSendMessage
	srv.handlers[models.EventEditMessage] = srv.handleEditMessage
	srv.handlers[models.EventDeleteMessage] = srv.handleDeleteMessage

	return srv
}

func (srv *EventService) HandleIncoming(
	ctx context.Context,
	roomID models.RoomId,
	userID models.UserId,
	data models.IncomingEvent,
) (models.ChatEvent, error) {
	if _, err := srv.roomStore.GetById(ctx, roomID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrForbidden
		}
		return nil, fmt.Errorf("ws get room id=%d: %w", roomID, err)
	}

	if _, err := srv.userStore.GetById(ctx, userID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrForbidden
		}
		return nil, fmt.Errorf("ws get user id=%d: %w", userID, err)
	}

	handler, ok := srv.handlers[models.IncomingEventType(data.Type)]
	if !ok {
		return nil, ErrUnsupportedEvent
	}

	return handler(ctx, roomID, userID, data)
}

func decodePayload(data models.IncomingEvent, v any) error {
	if err := json.Unmarshal(data.Data, v); err != nil {
		return ErrInvalidPayload
	}
	return nil
}

func (srv *EventService) ensureMember(
	ctx context.Context,
	roomID models.RoomId,
	userID models.UserId,
) error {
	exists, err := srv.roomMemberStore.Exists(ctx, roomID, userID)
	if err != nil {
		return fmt.Errorf("ws ensure member: %w", err)
	}
	if !exists {
		return ErrNotRoomMember
	}
	return nil
}

func (srv *EventService) handleJoinRoom(
	ctx context.Context,
	roomID models.RoomId,
	userID models.UserId,
	_ models.IncomingEvent,
) (models.ChatEvent, error) {
	return &models.BaseEvent{
		EventType: models.EventUserJoinedRoom,
		Data: map[string]any{
			"roomId": roomID,
			"userId": userID,
		},
	}, nil
}

func (srv *EventService) handleLeaveRoom(
	ctx context.Context,
	roomID models.RoomId,
	userID models.UserId,
	_ models.IncomingEvent,
) (models.ChatEvent, error) {
	return &models.BaseEvent{
		EventType: models.EventUserLeftRoom,
		Data: map[string]any{
			"roomId": roomID,
			"userId": userID,
		},
	}, nil
}

func (srv *EventService) handleSendMessage(
	ctx context.Context,
	roomID models.RoomId,
	userID models.UserId,
	data models.IncomingEvent,
) (models.ChatEvent, error) {
	if err := srv.ensureMember(ctx, roomID, userID); err != nil {
		return nil, err
	}

	var payload struct {
		Content string `json:"content"`
		Nonce   string `json:"nonce"`
	}

	if err := decodePayload(data, &payload); err != nil {
		return nil, err
	}

	if strings.TrimSpace(payload.Content) == "" {
		return nil, ErrInvalidPayload
	}

	msgID, err := srv.messageStore.Create(ctx, roomID, userID, payload.Content)
	if err != nil {
		return nil, fmt.Errorf("ws create message: %w", err)
	}

	msg, err := srv.messageStore.GetResponseById(ctx, msgID)
	if err != nil {
		return nil, fmt.Errorf("ws get response message: %w", err)
	}

	msg.Nonce = &payload.Nonce

	return &models.BaseEvent{
		EventType: models.EventMessageCreated,
		Data:      msg,
	}, nil
}

func (srv *EventService) handleEditMessage(
	ctx context.Context,
	roomID models.RoomId,
	userID models.UserId,
	data models.IncomingEvent,
) (models.ChatEvent, error) {
	if err := srv.ensureMember(ctx, roomID, userID); err != nil {
		return nil, err
	}

	var payload struct {
		MessageID models.MessageId `json:"messageId"`
		Content   string           `json:"content"`
	}

	if err := decodePayload(data, &payload); err != nil {
		return nil, err
	}

	if strings.TrimSpace(payload.Content) == "" {
		return nil, ErrInvalidPayload
	}

	msg, err := srv.messageStore.GetById(ctx, payload.MessageID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrForbidden
		}

		return nil, fmt.Errorf("edit ws get message=%d: %w", payload.MessageID, err)
	}

	if msg.UserId != userID {
		return nil, ErrForbidden
	}

	if err := srv.messageStore.UpdateContent(ctx, payload.MessageID, payload.Content); err != nil {
		return nil, fmt.Errorf("edit ws update content: %w", err)
	}

	updatedMsg, err := srv.messageStore.GetResponseById(ctx, payload.MessageID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrForbidden
		}

		return nil, fmt.Errorf("edit ws get response message: %w", err)
	}

	return &models.BaseEvent{
		EventType: models.EventMessageUpdated,
		Data:      updatedMsg,
	}, nil
}

func (srv *EventService) handleDeleteMessage(
	ctx context.Context,
	roomID models.RoomId,
	userID models.UserId,
	data models.IncomingEvent,
) (models.ChatEvent, error) {
	if err := srv.ensureMember(ctx, roomID, userID); err != nil {
		return nil, err
	}

	var payload struct {
		MessageID models.MessageId `json:"messageId"`
	}

	if err := decodePayload(data, &payload); err != nil {
		return nil, err
	}

	msg, err := srv.messageStore.GetById(ctx, payload.MessageID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrForbidden
		}

		return nil, fmt.Errorf("delete ws get message=%d: %w", payload.MessageID, err)
	}

	if msg.UserId != userID {
		return nil, ErrForbidden
	}

	if err := srv.messageStore.DeleteById(ctx, payload.MessageID); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, ErrForbidden
		}

		return nil, fmt.Errorf("delete ws delete message=%d: %w", payload.MessageID, err)
	}

	return &models.BaseEvent{
		EventType: models.EventMessageDeleted,
		Data: map[string]any{
			"messageId": payload.MessageID,
			"roomId":    roomID,
		},
	}, nil
}
