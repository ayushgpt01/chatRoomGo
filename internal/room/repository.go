package room

import "context"

type RoomStore interface {
	Create(ctx context.Context, name string) (RoomId, error)
	GetById(ctx context.Context, roomId RoomId) (*Room, error)
	UpdateName(ctx context.Context, roomId RoomId, name string) error
	Delete(ctx context.Context, roomId RoomId) error
}
