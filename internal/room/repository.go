package room

import (
	"context"

	"github.com/ayushgpt01/chatRoomGo/internal/user"
)

type RoomStore interface {
	Create(ctx context.Context, name string) (RoomId, error)
	GetById(ctx context.Context, roomId RoomId) (*Room, error)
	UpdateName(ctx context.Context, roomId RoomId, name string) error
	Delete(ctx context.Context, roomId RoomId) error
}

type RoomMemberStore interface {
	JoinRoom(ctx context.Context, roomId RoomId, userId user.UserId) error
	LeaveRoom(ctx context.Context, roomId RoomId, userId user.UserId) error
	Exists(ctx context.Context, roomId RoomId, userId user.UserId) (bool, error)
	CountByRoomId(ctx context.Context, roomId RoomId) (int, error)
	GetByRoomId(ctx context.Context, roomId RoomId) ([]user.UserId, error)
}
