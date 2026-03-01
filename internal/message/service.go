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
	return srv.roomMemberStore.Exists(ctx, roomID, userID)
	// if err != nil {
	// 	return err
	// }
	// if !exists {
	// 	return ErrNotRoomMember
	// }
	// return nil
}

func (srv *MessageService) HandleGetMessages(ctx context.Context, payload GetMessagesPayload) (*GetMessagesResponse, error) {
	exists, err := srv.ensureMember(ctx, payload.RoomId, payload.UserId)
	if err != nil {
		return &GetMessagesResponse{}, err
	}

	if !exists {
		return &GetMessagesResponse{}, fmt.Errorf("unauthorised user")
	}

	return srv.messageStore.GetMessagesById(ctx, payload.RoomId, payload.Limit, payload.Cursor)
}
