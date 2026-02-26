package auth

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
	_ "modernc.org/sqlite"
)

type AuthStore interface {
	SaveRefreshToken(ctx context.Context, userId models.UserId, token string, expiresAt time.Time) error
	ValidateRefreshToken(ctx context.Context, token string) (models.UserId, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	CleanupExpiredTokens(ctx context.Context) error
}

type SQLiteAuthRepo struct {
	db *sql.DB
}

func NewSQLiteAuthRepo(ctx context.Context, db *sql.DB) (*SQLiteAuthRepo, error) {
	store := SQLiteAuthRepo{db}

	if err := store.init(ctx); err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *SQLiteAuthRepo) init(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS refresh_tokens(
		user_id INTEGER NOT NULL,
		token TEXT NOT NULL,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		PRIMARY KEY (user_id, token)
	)`

	if _, err := s.db.ExecContext(ctx, createTableSQL); err != nil {
		return err
	}

	return nil
}

func (s *SQLiteAuthRepo) SaveRefreshToken(ctx context.Context, userId models.UserId, token string, expiresAt time.Time) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO refresh_tokens(user_id, token, expires_at) VALUES(?, ?, ?)", userId, token, expiresAt)
	if err != nil {
		return fmt.Errorf("SaveRefreshToken %d: %v", userId, err)
	}

	return nil
}

func (s *SQLiteAuthRepo) ValidateRefreshToken(ctx context.Context, token string) (models.UserId, error) {
	var userId models.UserId
	row := s.db.QueryRowContext(ctx, "SELECT user_id FROM refresh_tokens WHERE token = ? AND expires_at > CURRENT_TIMESTAMP", token)

	err := row.Scan(&userId)
	if err != nil {
		return 0, err
	}

	return userId, nil
}

func (s *SQLiteAuthRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE token = ?", token)
	if err != nil {
		return fmt.Errorf("DeleteRefreshToken %s: %v", token, err)
	}

	return nil
}

func (s *SQLiteAuthRepo) CleanupExpiredTokens(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM refresh_tokens WHERE expires_at < CURRENT_TIMESTAMP")
	return err
}
