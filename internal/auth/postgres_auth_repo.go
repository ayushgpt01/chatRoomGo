package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type PostgresAuthRepo struct {
	db *sql.DB
}

func NewPostgresAuthRepo(ctx context.Context, db *sql.DB) *PostgresAuthRepo {
	return &PostgresAuthRepo{db}
}

func (s *PostgresAuthRepo) SaveRefreshToken(ctx context.Context, userId models.UserId, token string, expiresAt time.Time) error {
	query := "INSERT INTO refresh_tokens(user_id, token, expires_at) VALUES($1, $2, $3)"

	_, err := s.db.ExecContext(ctx, query, userId, token, expiresAt)
	if err != nil {
		return fmt.Errorf("saving refresh token %d: %w", userId, err)
	}

	return nil
}

func (s *PostgresAuthRepo) ValidateRefreshToken(ctx context.Context, token string) (models.UserId, error) {
	var userId models.UserId
	query := "SELECT user_id FROM refresh_tokens WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP"
	row := s.db.QueryRowContext(ctx, query, token)

	err := row.Scan(&userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("getting refresh token %s: %w", token, models.ErrNotFound)
		}
		return 0, fmt.Errorf("scanning refresh token %s: %w", token, err)
	}

	return userId, nil
}

func (s *PostgresAuthRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE token = $1", token)
	if err != nil {
		return fmt.Errorf("deleting refresh token %s: %w", token, err)
	}

	return nil
}

func (s *PostgresAuthRepo) CleanupExpiredTokens(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE expires_at < CURRENT_TIMESTAMP")
	if err != nil {
		return fmt.Errorf("cleaning up expired tokens: %w", err)
	}

	return nil
}
