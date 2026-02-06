package chat

import (
	"context"

	"github.com/ayushgpt01/chatRoomGo/internal/dto"
	"github.com/ayushgpt01/chatRoomGo/internal/message"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
)

type ChatService struct {
	userStore       user.UserStore
	roomStore       room.RoomStore
	messageStore    message.MessageStore
	roomMemberStore RoomMemberStore
}

func NewChatService(
	userStore user.UserStore,
	roomStore room.RoomStore,
	messageStore message.MessageStore,
	roomMemberStore RoomMemberStore) *ChatService {
	return &ChatService{userStore, roomStore, messageStore, roomMemberStore}
}

func (srv *ChatService) HandleIncoming(
	ctx context.Context, roomId string, userId string, data dto.IncomingMessage) (ChatEvent, error) {
	// Check if room exists
	// Check if user can send message in this room
	// Parse message
	// Add it to db
	// Create an new_message event
	// Return it

	return &NewMessageEvent{}, nil
}
