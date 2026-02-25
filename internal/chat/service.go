package chat

import (
	"context"
	"encoding/json"

	"github.com/ayushgpt01/chatRoomGo/internal/message"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/types"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
)

type eventHandler func(
	ctx context.Context,
	roomID room.RoomId,
	userID user.UserId,
	data types.IncomingEvent,
) (types.ChatEvent, error)

type ChatService struct {
	userStore       user.UserStore
	roomStore       room.RoomStore
	messageStore    message.MessageStore
	roomMemberStore room.RoomMemberStore

	handlers map[types.IncomingEventType]eventHandler
}

func NewChatService(
	userStore user.UserStore,
	roomStore room.RoomStore,
	messageStore message.MessageStore,
	roomMemberStore room.RoomMemberStore,
) *ChatService {
	srv := &ChatService{
		userStore:       userStore,
		roomStore:       roomStore,
		messageStore:    messageStore,
		roomMemberStore: roomMemberStore,
		handlers:        make(map[types.IncomingEventType]eventHandler),
	}

	srv.handlers[types.EventJoinRoom] = srv.handleJoinRoom
	srv.handlers[types.EventLeaveRoom] = srv.handleLeaveRoom
	srv.handlers[types.EventSendMessage] = srv.handleSendMessage
	srv.handlers[types.EventEditMessage] = srv.handleEditMessage
	srv.handlers[types.EventDeleteMessage] = srv.handleDeleteMessage

	return srv
}

func (srv *ChatService) HandleIncoming(
	ctx context.Context,
	roomID room.RoomId,
	userID user.UserId,
	data types.IncomingEvent,
) (types.ChatEvent, error) {
	if _, err := srv.roomStore.GetById(ctx, roomID); err != nil {
		return nil, err
	}

	if _, err := srv.userStore.GetById(ctx, userID); err != nil {
		return nil, err
	}

	handler, ok := srv.handlers[types.IncomingEventType(data.Type)]
	if !ok {
		return nil, ErrUnsupportedEvent
	}

	return handler(ctx, roomID, userID, data)
}

func decodePayload(data types.IncomingEvent, v any) error {
	if err := json.Unmarshal(data.Data, v); err != nil {
		return ErrInvalidPayload
	}
	return nil
}

func (srv *ChatService) ensureMember(
	ctx context.Context,
	roomID room.RoomId,
	userID user.UserId,
) error {
	exists, err := srv.roomMemberStore.Exists(ctx, roomID, userID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrNotRoomMember
	}
	return nil
}

func (srv *ChatService) handleJoinRoom(
	ctx context.Context,
	roomID room.RoomId,
	userID user.UserId,
	_ types.IncomingEvent,
) (types.ChatEvent, error) {
	return &types.BaseEvent{
		EventType: types.EventUserJoinedRoom,
		Data: map[string]any{
			"roomId": roomID,
			"userId": userID,
		},
	}, nil
}

func (srv *ChatService) handleLeaveRoom(
	ctx context.Context,
	roomID room.RoomId,
	userID user.UserId,
	_ types.IncomingEvent,
) (types.ChatEvent, error) {
	return &types.BaseEvent{
		EventType: types.EventUserLeftRoom,
		Data: map[string]any{
			"roomId": roomID,
			"userId": userID,
		},
	}, nil
}

func (srv *ChatService) handleSendMessage(
	ctx context.Context,
	roomID room.RoomId,
	userID user.UserId,
	data types.IncomingEvent,
) (types.ChatEvent, error) {
	if err := srv.ensureMember(ctx, roomID, userID); err != nil {
		return nil, err
	}

	var payload struct {
		Content string `json:"content"`
	}

	if err := decodePayload(data, &payload); err != nil {
		return nil, err
	}

	msgID, err := srv.messageStore.Create(ctx, roomID, userID, payload.Content)
	if err != nil {
		return nil, err
	}

	msg, err := srv.messageStore.GetById(ctx, msgID)
	if err != nil {
		return nil, err
	}

	return &types.BaseEvent{
		EventType: types.EventMessageCreated,
		Data:      msg,
	}, nil
}

func (srv *ChatService) handleEditMessage(
	ctx context.Context,
	roomID room.RoomId,
	userID user.UserId,
	data types.IncomingEvent,
) (types.ChatEvent, error) {
	if err := srv.ensureMember(ctx, roomID, userID); err != nil {
		return nil, err
	}

	var payload struct {
		MessageID message.MessageId `json:"messageId"`
		Content   string            `json:"content"`
	}

	if err := decodePayload(data, &payload); err != nil {
		return nil, err
	}

	msg, err := srv.messageStore.GetById(ctx, payload.MessageID)
	if err != nil {
		return nil, err
	}

	if msg.UserId != userID {
		return nil, ErrForbidden
	}

	if err := srv.messageStore.UpdateContent(ctx, payload.MessageID, payload.Content); err != nil {
		return nil, err
	}

	updatedMsg, err := srv.messageStore.GetById(ctx, payload.MessageID)
	if err != nil {
		return nil, err
	}

	return &types.BaseEvent{
		EventType: types.EventMessageUpdated,
		Data:      updatedMsg,
	}, nil
}

func (srv *ChatService) handleDeleteMessage(
	ctx context.Context,
	roomID room.RoomId,
	userID user.UserId,
	data types.IncomingEvent,
) (types.ChatEvent, error) {
	if err := srv.ensureMember(ctx, roomID, userID); err != nil {
		return nil, err
	}

	var payload struct {
		MessageID message.MessageId `json:"messageId"`
	}

	if err := decodePayload(data, &payload); err != nil {
		return nil, err
	}

	msg, err := srv.messageStore.GetById(ctx, payload.MessageID)
	if err != nil {
		return nil, err
	}

	if msg.UserId != userID {
		return nil, ErrForbidden
	}

	if err := srv.messageStore.DeleteById(ctx, payload.MessageID); err != nil {
		return nil, err
	}

	return &types.BaseEvent{
		EventType: types.EventMessageDeleted,
		Data: map[string]any{
			"messageId": payload.MessageID,
			"roomId":    roomID,
		},
	}, nil
}
