package message

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type PostgresMessageRepo struct {
	db *sql.DB
}

func NewPostgresMessageRepo(ctx context.Context, db *sql.DB) *PostgresMessageRepo {
	return &PostgresMessageRepo{db}
}

func (s *PostgresMessageRepo) GetById(ctx context.Context, id models.MessageId) (*models.Message, error) {
	var message models.Message
	var updatedAt sql.NullTime

	query := `SELECT id, content, user_id, room_id, created_at, updated_at, delivered
	FROM messages
	WHERE id = $1 AND deleted_at IS NULL`

	row := s.db.QueryRowContext(ctx, query, id)

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

func (s *PostgresMessageRepo) Create(ctx context.Context, roomId models.RoomId, userId models.UserId, content string) (models.MessageId, error) {
	var messageId models.MessageId
	query := "INSERT INTO messages(user_id, room_id, content) VALUES($1, $2, $3) RETURNING id"

	err := s.db.QueryRowContext(ctx, query, userId, roomId, content).Scan(&messageId)
	if err != nil {
		return 0, fmt.Errorf("inserting message into room %d: %w", roomId, err)
	}

	return messageId, nil
}

func (s *PostgresMessageRepo) DeleteById(ctx context.Context, id models.MessageId) error {
	res, err := s.db.ExecContext(ctx, "UPDATE messages SET deleted_at = CURRENT_TIMESTAMP WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("soft deleting message by id %d: %w", id, err)
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

func (s *PostgresMessageRepo) UpdateContent(ctx context.Context, id models.MessageId, content string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE messages SET content = $1 WHERE id = $2", content, id)
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

func (s *PostgresMessageRepo) GetMessageReaders(ctx context.Context, messageId models.MessageId, roomId models.RoomId) ([]models.UserId, error) {
	query := `SELECT rm.user_id 
	FROM room_members rm 
	WHERE rm.room_id = $1 AND rm.last_message_read_id >= $2`

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

func (s *PostgresMessageRepo) GetResponseById(ctx context.Context, id models.MessageId) (*models.ResponseMessage, error) {
	var message models.ResponseMessage
	var updatedAt sql.NullTime
	var readByJSON []byte

	row := s.db.QueryRowContext(ctx, `SELECT m.id, m.content, m.updated_at, u.id, u.name, m.created_at, m.room_id, m.delivered,
	COALESCE((
		SELECT json_agg(rm.user_id) 
		FROM room_members rm 
		WHERE rm.room_id = m.room_id AND rm.last_message_read_id >= m.id
	), '[]'::json) as read_by
	FROM messages m
	JOIN users u ON m.user_id = u.id
	WHERE m.id = $1 AND m.deleted_at IS NULL`, id)

	err := row.Scan(
		&message.Id,
		&message.Content,
		&updatedAt,
		&message.SenderId,
		&message.SenderName,
		&message.SentAt,
		&message.RoomId,
		&message.Delivered,
		&readByJSON,
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

	if err := json.Unmarshal(readByJSON, &message.ReadBy); err != nil {
		message.ReadBy = []models.UserId{}
	}

	return &message, nil
}

func (s *PostgresMessageRepo) GetMessagesById(ctx context.Context, roomId models.RoomId, limit int, cursor *string) (*GetMessagesResponse, error) {
	query := `SELECT m.id, m.content, m.updated_at, u.id, u.name, m.created_at, m.room_id, m.delivered,
	COALESCE((
		SELECT json_agg(rm.user_id) 
		FROM room_members rm 
		WHERE rm.room_id = m.room_id AND rm.last_message_read_id >= m.id
	), '[]'::json) as read_by
	FROM messages m
	JOIN users u ON m.user_id = u.id
	WHERE m.room_id = $1 AND m.deleted_at IS NULL`

	args := []any{roomId}
	placeholderCount := 1

	if cursor != nil && *cursor != "" {
		placeholderCount++
		query += fmt.Sprintf(" AND m.id > $%d ", placeholderCount)
		args = append(args, *cursor)
	}

	placeholderCount++
	query += fmt.Sprintf(" ORDER BY m.id ASC LIMIT $%d", placeholderCount)
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
		var readByJSON []byte

		err := rows.Scan(
			&msg.Id,
			&msg.Content,
			&updatedAt,
			&msg.SenderId,
			&msg.SenderName,
			&msg.SentAt,
			&msg.RoomId,
			&msg.Delivered,
			&readByJSON,
		)

		if updatedAt.Valid {
			msg.EditedAt = &updatedAt.Time
		}

		if err != nil {
			return nil, fmt.Errorf("scanning messages for room %d: %w", roomId, err)
		}

		if err := json.Unmarshal(readByJSON, &msg.ReadBy); err != nil {
			msg.ReadBy = []models.UserId{}
		}

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

func (s *PostgresMessageRepo) MarkAsDelivered(ctx context.Context, messageId models.MessageId) error {
	res, err := s.db.ExecContext(ctx,
		"UPDATE messages SET delivered = TRUE WHERE id = $1",
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
