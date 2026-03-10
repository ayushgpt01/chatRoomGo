package message

import (
	"context"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
)

type MessageService struct {
	messageStore    MessageStore
	roomMemberStore room.RoomMemberStore
	hub             models.HubBroadcaster
}

func NewMessageService(messageStore MessageStore, roomMemberStore room.RoomMemberStore, hub models.HubBroadcaster) *MessageService {
	return &MessageService{messageStore, roomMemberStore, hub}
}

func (srv *MessageService) ensureMember(
	ctx context.Context,
	roomID models.RoomId,
	userID models.UserId,
) (bool, error) {
	exists, err := srv.roomMemberStore.Exists(ctx, roomID, userID)
	if err != nil {
		return false, fmt.Errorf("ensure member user_id=%d exists in room_id=%d: %w", userID, roomID, err)
	}

	return exists, nil
}

func (srv *MessageService) HandleGetMessages(ctx context.Context, payload GetMessagesPayload) (*GetMessagesResponse, error) {
	exists, err := srv.ensureMember(ctx, payload.RoomId, payload.UserId)
	if err != nil {
		return &GetMessagesResponse{}, fmt.Errorf("get messages: %w", err)
	}

	if !exists {
		return &GetMessagesResponse{}, models.ErrUnauthorized
	}

	response, err := srv.messageStore.GetMessagesById(ctx, payload.RoomId, payload.Limit, payload.Cursor)
	if err != nil {
		return &GetMessagesResponse{}, fmt.Errorf("get messages by room_id=%d: %w", payload.RoomId, err)
	}

	return response, nil
}

func (srv *MessageService) HandleSendMessage(
	ctx context.Context,
	payload SendMessagePayload,
) (*models.ResponseMessage, error) {

	exists, err := srv.ensureMember(ctx, payload.RoomId, payload.UserId)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}
	if !exists {
		return nil, models.ErrUnauthorized
	}

	id, err := srv.messageStore.Create(ctx, payload.RoomId, payload.UserId, payload.Content)
	if err != nil {
		return nil, fmt.Errorf("create message: %w", err)
	}

	res, err := srv.messageStore.GetResponseById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get response message id=%d: %w", id, err)
	}

	srv.hub.Broadcast(payload.RoomId, &models.MessageCreatedEvent{
		Data: models.MessageCreatedPayload{
			Message: res,
		},
	})

	return res, nil
}

func (srv *MessageService) HandleEditMessage(
	ctx context.Context,
	payload EditMessagePayload,
) (*models.ResponseMessage, error) {

	msg, err := srv.messageStore.GetById(ctx, payload.MessageId)
	if err != nil {
		return nil, fmt.Errorf("get message id=%d: %w", payload.MessageId, err)
	}

	if msg.UserId != payload.UserId {
		return nil, models.ErrUnauthorized
	}

	if msg.RoomId != payload.RoomId {
		return nil, models.ErrNotFound
	}

	err = srv.messageStore.UpdateContent(ctx, payload.MessageId, payload.Content)
	if err != nil {
		return nil, fmt.Errorf("update message id=%d: %w", payload.MessageId, err)
	}

	res, err := srv.messageStore.GetResponseById(ctx, payload.MessageId)
	if err != nil {
		return nil, fmt.Errorf("get response message id=%d: %w", payload.MessageId, err)
	}

	srv.hub.Broadcast(payload.RoomId, &models.MessageUpdatedEvent{
		Data: models.MessageUpdatedPayload{
			Message: res,
		},
	})

	return res, nil
}

func (srv *MessageService) HandleDeleteMessage(
	ctx context.Context,
	payload DeleteMessagePayload,
) error {
	msg, err := srv.messageStore.GetById(ctx, payload.MessageId)
	if err != nil {
		return fmt.Errorf("get message id=%d: %w", payload.MessageId, err)
	}

	if msg.UserId != payload.UserId {
		return models.ErrUnauthorized
	}

	if msg.RoomId != payload.RoomId {
		return models.ErrNotFound
	}

	err = srv.messageStore.DeleteById(ctx, payload.MessageId)
	if err != nil {
		return fmt.Errorf("delete message id=%d: %w", payload.MessageId, err)
	}

	srv.hub.Broadcast(payload.RoomId, &models.MessageDeletedEvent{
		Data: models.MessageDeletedPayload{
			MessageID: payload.MessageId,
			RoomID:    payload.RoomId,
		},
	})

	return nil
}

func (srv *MessageService) HandleMarkAsRead(
	ctx context.Context,
	messageId models.MessageId,
	userId models.UserId,
) error {
	msg, err := srv.messageStore.GetById(ctx, messageId)
	if err != nil {
		return fmt.Errorf("get message for read receipt id=%d: %w", messageId, err)
	}

	exists, err := srv.ensureMember(ctx, msg.RoomId, userId)
	if err != nil {
		return fmt.Errorf("check membership for read receipt: %w", err)
	}
	if !exists {
		return models.ErrUnauthorized
	}

	if err := srv.roomMemberStore.UpdateLastMessageRead(ctx, msg.RoomId, userId, messageId); err != nil {
		return fmt.Errorf("update last message read: %w", err)
	}

	return nil
}

func (srv *MessageService) HandleMarkAsDelivered(
	ctx context.Context,
	messageId models.MessageId,
) error {
	if err := srv.messageStore.MarkAsDelivered(ctx, messageId); err != nil {
		return fmt.Errorf("mark message as delivered id=%d: %w", messageId, err)
	}

	return nil
}
