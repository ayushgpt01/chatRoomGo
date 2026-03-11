package message

import (
	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type GetMessagesPayload struct {
	UserId models.UserId `json:"userId"`
	RoomId models.RoomId `json:"roomId"`
	Limit  int           `json:"limit"`
	Cursor *string       `json:"cursor"`
}

type GetMessagesResponse struct {
	Messages   []models.ResponseMessage `json:"messages"`
	NextCursor *string                  `json:"nextCursor"`
}

type SendMessagePayload struct {
	UserId  models.UserId
	RoomId  models.RoomId
	Content string
	Nonce   *string
}

type EditMessagePayload struct {
	UserId    models.UserId
	MessageId models.MessageId
	RoomId    models.RoomId
	Content   string
}

type DeleteMessagePayload struct {
	UserId    models.UserId
	MessageId models.MessageId
	RoomId    models.RoomId
}
