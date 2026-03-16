package user

import (
	"context"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type UserStore interface {
	GetById(ctx context.Context, id models.UserId) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, username string, name string, passwordHash string, role models.AccountRole) (models.UserId, error)
	UpdateName(ctx context.Context, id models.UserId, name string) error
	UpdateUsername(ctx context.Context, id models.UserId, username string) error
	DeleteById(ctx context.Context, id models.UserId) error
}
