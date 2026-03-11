package message

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
	_ "modernc.org/sqlite"
)

type MessageStore interface {
	GetById(ctx context.Context, id models.MessageId) (*models.Message, error)
	Create(ctx context.Context, roomId models.RoomId, userId models.UserId, content string) (models.MessageId, error)
	DeleteById(ctx context.Context, id models.MessageId) error
	UpdateContent(ctx context.Context, id models.MessageId, content string) error
	GetResponseById(ctx context.Context, id models.MessageId) (*models.ResponseMessage, error)
	GetMessagesById(ctx context.Context, roomId models.RoomId, limit int, cursor *string) (*GetMessagesResponse, error)
	MarkAsDelivered(ctx context.Context, messageId models.MessageId) error
}

type SQLiteMessageRepo struct {
	db *sql.DB
}

func NewSQLiteMessageRepo(ctx context.Context, db *sql.DB) (*SQLiteMessageRepo, error) {
	store := SQLiteMessageRepo{db}

	if err := store.init(ctx); err != nil {
		return nil, fmt.Errorf("initializing messages table: %w", err)
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
		updated_at DATETIME DEFAULT NULL,
		delivered BOOLEAN DEFAULT TRUE,
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
		return fmt.Errorf("creating messages table: %w", err)
	}

	if _, err := s.db.ExecContext(ctx, createTriggerSQL); err != nil {
		return fmt.Errorf("creating update_message_timestamp trigger: %w", err)
	}

	if _, err := s.db.ExecContext(ctx, createUserIdIndexSQL); err != nil {
		return fmt.Errorf("creating messages user_id index: %w", err)
	}

	if _, err := s.db.ExecContext(ctx, createRoomIdIndexSQL); err != nil {
		return fmt.Errorf("creating messages room_id index: %w", err)
	}

	return nil
}

func (s *SQLiteMessageRepo) GetById(ctx context.Context, id models.MessageId) (*models.Message, error) {
	var message models.Message
	var updatedAt sql.NullTime

	row := s.db.QueryRowContext(ctx, `SELECT id, content, user_id, room_id, created_at, updated_at, delivered
	FROM messages
	WHERE id = ?`, id)

	err := row.Scan(
		&message.Id,
		&message.Content,
		&message.UserId,
		&message.RoomId,
		&message.CreatedAt,
		&updatedAt,
		&message.Delivered,
	)

	if updatedAt.Valid {
		message.UpdatedAt = &updatedAt.Time
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting message by id %d: %w", id, models.ErrNotFound)
		}
		return nil, fmt.Errorf("scanning message by id %d: %w", id, err)
	}

	return &message, nil
}

func (s *SQLiteMessageRepo) Create(ctx context.Context, roomId models.RoomId, userId models.UserId, content string) (models.MessageId, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO messages(user_id, room_id, content) VALUES(?, ?, ?)", userId, roomId, content)
	if err != nil {
		return 0, fmt.Errorf("inserting message into room %d: %w", roomId, err)
	}

	messageId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("getting last insert id for message: %w", err)
	}

	return messageId, nil
}

func (s *SQLiteMessageRepo) DeleteById(ctx context.Context, id models.MessageId) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM messages WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("deleting message by id %d: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected for delete message %d: %w", id, err)
	}

	if count == 0 {
		return fmt.Errorf("deleting message by id %d: %w", id, models.ErrNotFound)
	}

	return nil
}

func (s *SQLiteMessageRepo) UpdateContent(ctx context.Context, id models.MessageId, content string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE messages SET content = ? WHERE id = ?", content, id)
	if err != nil {
		return fmt.Errorf("updating message content for id %d: %w", id, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected for update message %d: %w", id, err)
	}

	if count == 0 {
		return fmt.Errorf("updating message content for id %d: %w", id, models.ErrNotFound)
	}

	return nil
}

func (s *SQLiteMessageRepo) GetMessageReaders(ctx context.Context, messageId models.MessageId, roomId models.RoomId) ([]models.UserId, error) {
	query := `SELECT rm.user_id 
	FROM room_members rm 
	WHERE rm.room_id = ? AND rm.last_message_read_id >= ?`

	rows, err := s.db.QueryContext(ctx, query, roomId, messageId)
	if err != nil {
		return nil, fmt.Errorf("getting message readers for message %d: %w", messageId, err)
	}
	defer rows.Close()

	var readers []models.UserId
	for rows.Next() {
		var userId models.UserId
		if err := rows.Scan(&userId); err != nil {
			return nil, fmt.Errorf("scanning message reader: %w", err)
		}
		readers = append(readers, userId)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating message readers: %w", err)
	}

	return readers, nil
}

func (s *SQLiteMessageRepo) GetResponseById(ctx context.Context, id models.MessageId) (*models.ResponseMessage, error) {
	var message models.ResponseMessage
	var updatedAt sql.NullTime

	row := s.db.QueryRowContext(ctx, `SELECT m.id, m.content, m.updated_at, u.id, u.name, m.created_at, m.room_id, m.delivered
	FROM messages m
	JOIN users u ON m.user_id = u.id
	WHERE m.id = ?`, id)

	err := row.Scan(
		&message.Id,
		&message.Content,
		&updatedAt,
		&message.SenderId,
		&message.SenderName,
		&message.SentAt,
		&message.RoomId,
		&message.Delivered,
	)

	if updatedAt.Valid {
		message.EditedAt = &updatedAt.Time
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting response message by id %d: %w", id, models.ErrNotFound)
		}
		return nil, fmt.Errorf("scanning response message by id %d: %w", id, err)
	}

	// Get readers for this message
	readers, err := s.GetMessageReaders(ctx, id, message.RoomId)
	if err != nil {
		// Don't fail the entire operation if we can't get readers
		readers = []models.UserId{}
	}

	message.ReadBy = readers

	return &message, nil
}

func (s *SQLiteMessageRepo) GetMessagesById(ctx context.Context, roomId models.RoomId, limit int, cursor *string) (*GetMessagesResponse, error) {
	query := `SELECT m.id, m.content, m.updated_at, u.id, u.name, m.created_at, m.room_id, m.delivered
	FROM messages m
	JOIN users u ON m.user_id = u.id
	WHERE m.room_id = ?`

	args := []any{roomId}

	if cursor != nil && *cursor != "" {
		query += " AND m.id > ? "
		args = append(args, *cursor)
	}

	query += " ORDER BY m.id ASC LIMIT ?"
	args = append(args, limit+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("querying messages for room %d: %w", roomId, err)
	}

	defer rows.Close()

	messages := []models.ResponseMessage{}
	for rows.Next() {
		var msg models.ResponseMessage
		var updatedAt sql.NullTime

		err := rows.Scan(
			&msg.Id,
			&msg.Content,
			&updatedAt,
			&msg.SenderId,
			&msg.SenderName,
			&msg.SentAt,
			&msg.RoomId,
			&msg.Delivered,
		)

		if updatedAt.Valid {
			msg.EditedAt = &updatedAt.Time
		}

		if err != nil {
			return nil, fmt.Errorf("scanning messages for room %d: %w", roomId, err)
		}

		// Get readers for this message
		readers, err := s.GetMessageReaders(ctx, msg.Id, roomId)
		if err != nil {
			// Don't fail the entire operation if we can't get readers
			readers = []models.UserId{}
		}

		msg.ReadBy = readers
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating messages for room %d: %w", roomId, err)
	}

	var nextCursor *string
	if len(messages) > limit {
		lastRoom := messages[limit]
		c := fmt.Sprintf("%d", lastRoom.Id)
		nextCursor = &c
		messages = messages[:limit]
	}

	return &GetMessagesResponse{Messages: messages, NextCursor: nextCursor}, nil
}

func (s *SQLiteMessageRepo) MarkAsDelivered(ctx context.Context, messageId models.MessageId) error {
	res, err := s.db.ExecContext(ctx,
		"UPDATE messages SET delivered = TRUE WHERE id = ?",
		messageId)
	if err != nil {
		return fmt.Errorf("marking message %d as delivered: %w", messageId, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected for mark as delivered %d: %w", messageId, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("marking message %d as delivered: %w", messageId, models.ErrNotFound)
	}

	return nil
}
