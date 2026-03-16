package message

import (
	"context"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type MessageStore interface {
	GetById(ctx context.Context, id models.MessageId) (*models.Message, error)
	Create(ctx context.Context, roomId models.RoomId, userId models.UserId, content string) (models.MessageId, error)
	DeleteById(ctx context.Context, id models.MessageId) error
	UpdateContent(ctx context.Context, id models.MessageId, content string) error
	GetResponseById(ctx context.Context, id models.MessageId) (*models.ResponseMessage, error)
	GetMessagesById(ctx context.Context, roomId models.RoomId, limit int, cursor *string) (*GetMessagesResponse, error)
	MarkAsDelivered(ctx context.Context, messageId models.MessageId) error
}
