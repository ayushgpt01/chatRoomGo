package room

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type PostgresRoomMemberRepo struct {
	db *sql.DB
}

func NewPostgresRoomMemberRepo(ctx context.Context, db *sql.DB) *PostgresRoomMemberRepo {
	return &PostgresRoomMemberRepo{db}
}

func (s *PostgresRoomMemberRepo) JoinRoom(ctx context.Context, roomId models.RoomId, userId models.UserId) error {
	query := "INSERT INTO room_members(room_id, user_id) VALUES($1, $2) ON CONFLICT DO NOTHING"

	_, err := s.db.ExecContext(ctx, query, roomId, userId)
	if err != nil {
		return fmt.Errorf("join room room_id=%d user_id=%d: %w", roomId, userId, err)
	}

	return nil
}

func (s *PostgresRoomMemberRepo) LeaveRoom(ctx context.Context, roomId models.RoomId, userId models.UserId) error {
	query := "DELETE FROM room_members WHERE room_id = $1 AND user_id = $2"

	res, err := s.db.ExecContext(ctx, query, roomId, userId)
	if err != nil {
		return fmt.Errorf("leave room room_id=%d user_id=%d: %w", roomId, userId, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("leave room rows affected room_id=%d user_id=%d: %w", roomId, userId, err)
	}

	if count == 0 {
		return fmt.Errorf("leave room room_id=%d user_id=%d: %w", roomId, userId, models.ErrNotFound)
	}

	return nil
}

func (s *PostgresRoomMemberRepo) Exists(ctx context.Context, roomId models.RoomId, userId models.UserId) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = $1 AND user_id = $2)"

	var exists bool

	err := s.db.QueryRowContext(ctx, query, roomId, userId).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check room member exists room_id=%d user_id=%d: %w", roomId, userId, err)
	}

	return exists, nil
}

func (s *PostgresRoomMemberRepo) CountByRoomId(ctx context.Context, roomId models.RoomId) (int, error) {
	query := "SELECT COUNT(user_id) FROM room_members WHERE room_id = $1"
	var count int

	err := s.db.QueryRowContext(ctx, query, roomId).Scan(&count)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, fmt.Errorf("count room members room_id=%d: %w", roomId, err)
	}

	return count, nil
}

func (s *PostgresRoomMemberRepo) GetByRoomId(ctx context.Context, roomId models.RoomId) ([]models.UserId, error) {
	query := "SELECT user_id FROM room_members WHERE room_id = $1"

	rows, err := s.db.QueryContext(ctx, query, roomId)
	if err != nil {
		return nil, fmt.Errorf("get room members room_id=%d: %w", roomId, err)
	}
	defer rows.Close()

	var ids []models.UserId
	for rows.Next() {
		var id models.UserId
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan room member room_id=%d: %w", roomId, err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate room members room_id=%d: %w", roomId, err)
	}

	return ids, nil
}

func (s *PostgresRoomMemberRepo) GetRoomsByUserId(ctx context.Context, userId models.UserId, limit int, cursor *string) ([]*models.Room, *string, error) {
	query := `SELECT r.id, r.name, r.created_at, r.updated_at
	FROM rooms r
	JOIN room_members rm ON r.id = rm.room_id
	WHERE rm.user_id = $1`

	args := []any{userId}
	placeholderCount := 1

	if cursor != nil && *cursor != "" {
		placeholderCount++
		query += fmt.Sprintf(" AND r.id < $%d ", placeholderCount)
		args = append(args, *cursor)
	}

	placeholderCount++
	query += fmt.Sprintf(" ORDER BY r.id DESC LIMIT $%d", placeholderCount)
	args = append(args, limit+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, fmt.Errorf("querying rooms for user %d: %w", userId, err)
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(&room.Id, &room.Name, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, nil, fmt.Errorf("scanning rooms for user %d: %w", userId, err)
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterating rooms for user %d: %w", userId, err)
	}

	var nextCursor *string
	if len(rooms) > limit {
		lastRoom := rooms[limit-1]
		c := fmt.Sprintf("%d", lastRoom.Id)
		nextCursor = &c
		rooms = rooms[:limit]
	}

	return rooms, nextCursor, nil
}

func (s *PostgresRoomMemberRepo) UpdateLastMessageRead(ctx context.Context, roomId models.RoomId, userId models.UserId, messageId models.MessageId) error {
	res, err := s.db.ExecContext(ctx,
		"UPDATE room_members SET last_message_read_id = $1 WHERE room_id = $2 AND user_id = $3",
		messageId, roomId, userId)
	if err != nil {
		return fmt.Errorf("update last message read room_id=%d user_id=%d: %w", roomId, userId, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected for update last message read %d: %w", roomId, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("update last message read room_id=%d user_id=%d: %w", roomId, userId, models.ErrNotFound)
	}

	return nil
}

func (s *PostgresRoomMemberRepo) GetLastMessageRead(ctx context.Context, roomId models.RoomId, userId models.UserId) (models.MessageId, error) {
	query := "SELECT last_message_read_id FROM room_members WHERE room_id = $1 AND user_id = $2"
	var lastReadId models.MessageId

	err := s.db.QueryRowContext(ctx, query, roomId, userId).Scan(&lastReadId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("get last message read room_id=%d user_id=%d: %w", roomId, userId, models.ErrNotFound)
		}
		return 0, fmt.Errorf("get last message read room_id=%d user_id=%d: %w", roomId, userId, err)
	}

	return lastReadId, nil
}

func (s *PostgresRoomMemberRepo) GetRoomMembers(ctx context.Context, roomId models.RoomId) ([]*models.User, error) {
	query := `SELECT u.id, u.name, u.user_name, u.account_role, u.created_at, u.updated_at
	FROM users u
	JOIN room_members rm ON u.id = rm.user_id
	WHERE rm.room_id = $1
	ORDER BY u.created_at ASC
	LIMIT 50`

	rows, err := s.db.QueryContext(ctx, query, roomId)
	if err != nil {
		return nil, fmt.Errorf("get room members room_id=%d: %w", roomId, err)
	}
	defer rows.Close()

	var members []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(&user.Id, &user.Name, &user.Username, &user.AccountRole, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan room member room_id=%d: %w", roomId, err)
		}
		members = append(members, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate room members room_id=%d: %w", roomId, err)
	}

	return members, nil
}
