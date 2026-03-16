package room

import (
	"context"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type RoomStore interface {
	Create(ctx context.Context, name string) (*models.Room, error)
	GetById(ctx context.Context, roomId models.RoomId) (*models.Room, error)
	UpdateName(ctx context.Context, roomId models.RoomId, name string) error
	Delete(ctx context.Context, roomId models.RoomId) error
}

type RoomMemberStore interface {
	JoinRoom(ctx context.Context, roomId models.RoomId, userId models.UserId) error
	LeaveRoom(ctx context.Context, roomId models.RoomId, userId models.UserId) error
	Exists(ctx context.Context, roomId models.RoomId, userId models.UserId) (bool, error)
	CountByRoomId(ctx context.Context, roomId models.RoomId) (int, error)
	GetByRoomId(ctx context.Context, roomId models.RoomId) ([]models.UserId, error)
	GetRoomsByUserId(ctx context.Context, userId models.UserId, limit int, cursor *string) ([]*models.Room, *string, error)
	UpdateLastMessageRead(ctx context.Context, roomId models.RoomId, userId models.UserId, messageId models.MessageId) error
	GetLastMessageRead(ctx context.Context, roomId models.RoomId, userId models.UserId) (models.MessageId, error)
	GetRoomMembers(ctx context.Context, roomId models.RoomId) ([]*models.User, error)
}
