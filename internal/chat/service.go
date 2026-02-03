package chat

import (
	"github.com/ayushgpt01/chatRoomGo/internal/message"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
)

type ChatService struct {
	userStore    user.UserStore
	messageStore message.MessageStore
}

func NewChatService(userStore user.UserStore, messageStore message.MessageStore) *ChatService {
	return &ChatService{userStore, messageStore}
}

func (srv *ChatService) HandleIncoming(roomId string, userId string, data []byte) {

}
