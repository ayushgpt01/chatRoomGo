package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type PostgresUserRepo struct {
	db *sql.DB
}

func NewPostgresUserRepo(ctx context.Context, db *sql.DB) *PostgresUserRepo {
	return &PostgresUserRepo{db}
}

func (s *PostgresUserRepo) GetById(ctx context.Context, id models.UserId) (*models.User, error) {
	var user models.User

	row := s.db.QueryRowContext(ctx, `SELECT id, name, user_name, created_at, updated_at, password_hash, account_role 
	FROM users 
	WHERE id = $1`, id)

	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.Password, &user.AccountRole)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user by id=%d: %w", id, models.ErrNotFound)
		}

		return nil, fmt.Errorf("get user by id=%d: %w", id, err)
	}

	return &user, nil
}

func (s *PostgresUserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	row := s.db.QueryRowContext(ctx, `SELECT id, name, user_name, created_at, updated_at, password_hash, account_role 
	FROM users 
	WHERE user_name = $1`, username)

	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.Password, &user.AccountRole)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user by username=%s: %w", username, models.ErrNotFound)
		}
		return nil, fmt.Errorf("get user by username=%s: %w", username, err)
	}

	return &user, nil
}

func (s *PostgresUserRepo) Create(ctx context.Context, username, name, passwordHash string, role models.AccountRole) (models.UserId, error) {
	if role == "" {
		role = models.AccountRoleUser
	}

	var userId models.UserId
	query := "INSERT INTO users(name, user_name, password_hash, account_role) VALUES($1, $2, $3, $4) RETURNING id"

	err := s.db.QueryRowContext(ctx, query, name, username, passwordHash, role).Scan(&userId)
	if err != nil {
		return 0, fmt.Errorf("create user username=%s: %w", username, err)
	}

	return userId, nil
}

func (s *PostgresUserRepo) UpdateName(ctx context.Context, id models.UserId, name string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE users SET name = $1 WHERE id = $2", name, id)
	if err != nil {
		return fmt.Errorf("update name id=%d: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update name rows affected id=%d: %w", id, err)
	}

	if count == 0 {
		return fmt.Errorf("update name id=%d: %w", id, models.ErrNotFound)
	}

	return nil
}

func (s *PostgresUserRepo) UpdateUsername(ctx context.Context, id models.UserId, username string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE users SET user_name = $1 WHERE id = $2", username, id)
	if err != nil {
		return fmt.Errorf("update user name id=%d: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update user name rows affected id=%d: %w", id, err)
	}

	if count == 0 {
		return fmt.Errorf("update user name id=%d: %w", id, models.ErrNotFound)
	}

	return nil
}

func (s *PostgresUserRepo) DeleteById(ctx context.Context, id models.UserId) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete user id=%d: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete user rows affected id=%d: %w", id, err)
	}

	if count == 0 {
		return fmt.Errorf("delete user id=%d: %w", id, models.ErrNotFound)
	}

	return nil
}
