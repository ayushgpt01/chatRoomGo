package message

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
)

type SQLiteMessageRepo struct {
	db *sql.DB
}

func NewSQLiteMessageRepo(ctx context.Context, db *sql.DB) (*SQLiteMessageRepo, error) {
	store := SQLiteMessageRepo{db}

	if err := store.init(ctx); err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *SQLiteMessageRepo) init(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT,
		user_id INTEGER NOT NULL,
		room_id INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE RESTRICT,
		FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE RESTRICT
	);`

	createTriggerSQL := `CREATE TRIGGER IF NOT EXISTS update_message_timestamp
	AFTER UPDATE ON messages
	BEGIN
		UPDATE messages SET updated_at = CURRENT_TIMESTAMP WHERE ID = old.id;
	END;`

	createUserIdIndexSQL := `CREATE INDEX IF NOT EXISTS idx_messages_users_id ON messages(user_id)`
	createRoomIdIndexSQL := `CREATE INDEX IF NOT EXISTS idx_messages_rooms_id ON messages(room_id)`

	if _, err := s.db.ExecContext(ctx, createTableSQL); err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, createTriggerSQL); err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, createUserIdIndexSQL); err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, createRoomIdIndexSQL); err != nil {
		return err
	}

	return nil
}

func (s *SQLiteMessageRepo) GetById(ctx context.Context, id MessageId) (*Message, error) {
	var message Message

	row := s.db.QueryRowContext(ctx, `SELECT id, content, user_id, room_id, created_at, updated_at
	FROM messages
	WHERE id = ?`, id)

	err := row.Scan(&message.Id, &message.Content, &message.UserId, &message.RoomId, &message.CreatedAt, &message.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("GetById %d: no such message", id)
		}

		return nil, fmt.Errorf("GetById %d: %v", id, err)
	}

	return &message, nil
}

func (s *SQLiteMessageRepo) Create(ctx context.Context, roomId room.RoomId, userId user.UserId, content string) (MessageId, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO messages(user_id, room_id, content) VALUES(?, ?, ?)", userId, roomId, content)
	if err != nil {
		return 0, err
	}

	messageId, err := res.LastInsertId()
	return messageId, err
}

func (s *SQLiteMessageRepo) DeleteById(ctx context.Context, id MessageId) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Message not found for id: %d", id)
	}

	return nil
}

func (s *SQLiteMessageRepo) UpdateContent(ctx context.Context, id MessageId, content string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE messages SET content = ? WHERE id = ?", content, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Message not found for id: %d", id)
	}

	return nil
}
