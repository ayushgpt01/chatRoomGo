package message

import (
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type ResponseMessage struct {
	Id         models.MessageId `json:"id"`
	Content    string           `json:"content"`
	SenderName string           `json:"senderName"`
	SenderId   models.UserId    `json:"senderId"`
	SentAt     time.Time        `json:"sentAt"`
	EditedAt   *time.Time       `json:"editedAt"`
	Read       bool             `json:"read"`
	Nonce      *string          `json:"nonce,omitempty"`
}

type GetMessagesPayload struct {
	UserId models.UserId `json:"userId"`
	RoomId models.RoomId `json:"roomId"`
	Limit  int           `json:"limit"`
	Cursor *string       `json:"cursor"`
}

type GetMessagesResponse struct {
	Messages   []ResponseMessage `json:"messages"`
	NextCursor *string           `json:"nextCursor"`
}
