package user

import "context"

type UserStore interface {
	GetById(ctx context.Context, id UserId) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, username string, name string) (UserId, error)
	UpdateName(ctx context.Context, id UserId, name string) error
	UpdateUsername(ctx context.Context, id UserId, username string) error
	DeleteById(ctx context.Context, id UserId) error
}
