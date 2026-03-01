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
}

func NewMessageService(messageStore MessageStore, roomMemberStore room.RoomMemberStore) *MessageService {
	return &MessageService{messageStore, roomMemberStore}
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
