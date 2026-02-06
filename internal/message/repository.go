package message

import (
	"context"

	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
)

type MessageStore interface {
	GetById(ctx context.Context, id MessageId) (*Message, error)
	Create(ctx context.Context, roomId room.RoomId, userId user.UserId, content string) (MessageId, error)
	DeleteById(ctx context.Context, id MessageId) error
	UpdateContent(ctx context.Context, id MessageId, content string) error
}
