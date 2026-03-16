package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
	_ "modernc.org/sqlite"
)

type SQLiteUserRepo struct {
	db *sql.DB
}

func NewSqliteUserRepo(ctx context.Context, db *sql.DB) (*SQLiteUserRepo, error) {
	store := SQLiteUserRepo{db}

	if err := store.init(ctx); err != nil {
		return nil, fmt.Errorf("init user repo: %w", err)
	}

	return &store, nil
}

func (s *SQLiteUserRepo) init(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		user_name VARCHAR(255) NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		account_role TEXT NOT NULL DEFAULT 'user'
			CHECK(account_role IN ('user', 'admin', 'guest')),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createTriggerSQL := `CREATE TRIGGER IF NOT EXISTS update_user_timestamp
	AFTER UPDATE ON users
	BEGIN
		UPDATE users SET updated_at = CURRENT_TIMESTAMP where id = old.id;
	END;
	`

	if _, err := s.db.ExecContext(ctx, createTableSQL); err != nil {
		return fmt.Errorf("create users table: %w", err)
	}

	if _, err := s.db.ExecContext(ctx, createTriggerSQL); err != nil {
		return fmt.Errorf("create user update trigger: %w", err)
	}

	return nil
}

func (s *SQLiteUserRepo) GetById(ctx context.Context, id models.UserId) (*models.User, error) {
	var user models.User

	row := s.db.QueryRowContext(ctx, `SELECT id, name, user_name, created_at, updated_at, password_hash, account_role 
	FROM users 
	WHERE id = ?`, id)

	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.Password, &user.AccountRole)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user by id=%d: %w", id, models.ErrNotFound)
		}

		return nil, fmt.Errorf("get user by id=%d: %w", id, err)
	}

	return &user, nil
}

func (s *SQLiteUserRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	row := s.db.QueryRowContext(ctx, `SELECT id, name, user_name, created_at, updated_at, password_hash, account_role 
	FROM users 
	WHERE user_name = ?`, username)

	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.CreatedAt, &user.UpdatedAt, &user.Password, &user.AccountRole)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get user by username=%s: %w", username, models.ErrNotFound)
		}
		return nil, fmt.Errorf("get user by username=%s: %w", username, err)
	}

	return &user, nil
}

func (s *SQLiteUserRepo) Create(ctx context.Context, username, name, passwordHash string, role models.AccountRole) (models.UserId, error) {
	if role == "" {
		role = models.AccountRoleUser
	}

	query := "INSERT INTO users(name, user_name, password_hash, account_role) VALUES(?, ?, ?, ?)"
	res, err := s.db.ExecContext(ctx, query, name, username, passwordHash, role)
	if err != nil {
		return 0, fmt.Errorf("create user username=%s: %w", username, err)
	}

	userId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get last insert id username=%s: %w", username, err)
	}

	return userId, nil
}

func (s *SQLiteUserRepo) UpdateName(ctx context.Context, id models.UserId, name string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE users SET name = ? WHERE id = ?", name, id)
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

func (s *SQLiteUserRepo) UpdateUsername(ctx context.Context, id models.UserId, username string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE users SET user_name = ? WHERE id = ?", username, id)
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

func (s *SQLiteUserRepo) DeleteById(ctx context.Context, id models.UserId) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
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
