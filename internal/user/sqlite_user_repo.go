package user

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type SQLiteUserRepo struct {
	db *sql.DB
}

func NewSqliteUserRepo(ctx context.Context, db *sql.DB) (*SQLiteUserRepo, error) {
	store := SQLiteUserRepo{db}

	if err := store.init(ctx); err != nil {
		return nil, fmt.Errorf("Could not initalise store: %v", err)
	}

	return &store, nil
}

func (s *SQLiteUserRepo) init(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		user_name VARCHAR(255) NOT NULL UNIQUE,
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
		return err
	}

	if _, err := s.db.ExecContext(ctx, createTriggerSQL); err != nil {
		return err
	}

	return nil
}

func (s *SQLiteUserRepo) GetById(ctx context.Context, id UserId) (*User, error) {
	var user User

	row := s.db.QueryRowContext(ctx, `SELECT id, name, user_name, created_at, updated_at 
	FROM users 
	WHERE id = ?`, id)

	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return &user, fmt.Errorf("GetById %d: no such user", id)
		}

		return &user, fmt.Errorf("GetById %d: %v", id, err)
	}

	return &user, nil
}

func (s *SQLiteUserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	var user User

	row := s.db.QueryRowContext(ctx, `SELECT id, name, user_name, created_at, updated_at 
	FROM users 
	WHERE user_name = ?`, username)

	err := row.Scan(&user.Id, &user.Name, &user.Username, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return &user, fmt.Errorf("GetByUsername %s: no such user", username)
		}

		return &user, fmt.Errorf("GetByUsername %s: %v", username, err)
	}

	return &user, nil
}

func (s *SQLiteUserRepo) Create(ctx context.Context, username string, name string) (UserId, error) {
	var userId UserId

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return userId, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO users(name, user_name) VALUES(?, ?)")
	if err != nil {
		return userId, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, name, username)
	if err != nil {
		return userId, err
	}

	err = tx.Commit()
	if err != nil {
		return userId, err
	}

	userId, err = res.LastInsertId()
	return userId, err
}

func (s *SQLiteUserRepo) UpdateName(ctx context.Context, id UserId, name string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE users SET name = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, name, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("User not found for id: %d", id)
	}

	return nil
}

func (s *SQLiteUserRepo) UpdateUsername(ctx context.Context, id UserId, username string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE users SET user_name = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, username, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("User not found for id: %d", id)
	}

	return nil
}

func (s *SQLiteUserRepo) DeleteById(ctx context.Context, id UserId) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("User not found for id: %d", id)
	}

	return nil
}
