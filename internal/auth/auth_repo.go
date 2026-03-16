package auth

import (
	"context"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type AuthStore interface {
	SaveRefreshToken(ctx context.Context, userId models.UserId, token string, expiresAt time.Time) error
	ValidateRefreshToken(ctx context.Context, token string) (models.UserId, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	CleanupExpiredTokens(ctx context.Context) error
}
