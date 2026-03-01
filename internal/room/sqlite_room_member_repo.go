package room

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
	_ "modernc.org/sqlite"
)

type RoomMemberStore interface {
	JoinRoom(ctx context.Context, roomId models.RoomId, userId models.UserId) error
	LeaveRoom(ctx context.Context, roomId models.RoomId, userId models.UserId) error
	Exists(ctx context.Context, roomId models.RoomId, userId models.UserId) (bool, error)
	CountByRoomId(ctx context.Context, roomId models.RoomId) (int, error)
	GetByRoomId(ctx context.Context, roomId models.RoomId) ([]models.UserId, error)
	GetRoomsByUserId(ctx context.Context, userId models.UserId, limit int, cursor *string) ([]*models.Room, *string, error)
}

type SQLiteRoomMemberRepo struct {
	db *sql.DB
}

func NewSQLiteRoomMemberRepo(ctx context.Context, db *sql.DB) (*SQLiteRoomMemberRepo, error) {
	store := SQLiteRoomMemberRepo{db}

	if err := store.init(ctx); err != nil {
		return nil, err
	}

	return &store, nil
}

func (s *SQLiteRoomMemberRepo) init(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS room_members(
		room_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
		PRIMARY KEY (room_id, user_id)
	)`

	createUserIdIndexSQL := `CREATE INDEX IF NOT EXISTS idx_room_members_users_id ON room_members(user_id)`

	if _, err := s.db.ExecContext(ctx, createTableSQL); err != nil {
		return err
	}

	if _, err := s.db.ExecContext(ctx, createUserIdIndexSQL); err != nil {
		return err
	}

	return nil
}

func (s *SQLiteRoomMemberRepo) JoinRoom(ctx context.Context, roomId models.RoomId, userId models.UserId) error {
	_, err := s.db.ExecContext(ctx, "INSERT OR IGNORE INTO room_members(room_id, user_id) VALUES(?, ?)", roomId, userId)
	return err
}

func (s *SQLiteRoomMemberRepo) LeaveRoom(ctx context.Context, roomId models.RoomId, userId models.UserId) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM room_members WHERE room_id = ? AND user_id = ?", roomId, userId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("User %d has not joined the room: %d", userId, roomId)
	}

	return nil
}

func (s *SQLiteRoomMemberRepo) Exists(ctx context.Context, roomId models.RoomId, userId models.UserId) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM room_members WHERE room_id = ? AND user_id = ?)"

	var exists bool

	err := s.db.QueryRowContext(ctx, query, roomId, userId).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

func (s *SQLiteRoomMemberRepo) CountByRoomId(ctx context.Context, roomId models.RoomId) (int, error) {
	query := "SELECT COUNT(user_id) FROM room_members WHERE room_id = ?"
	var count int

	err := s.db.QueryRowContext(ctx, query, roomId).Scan(&count)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	return count, nil
}

func (s *SQLiteRoomMemberRepo) GetByRoomId(ctx context.Context, roomId models.RoomId) ([]models.UserId, error) {
	query := "SELECT user_id FROM room_members WHERE room_id = ?"

	rows, err := s.db.QueryContext(ctx, query, roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []models.UserId
	for rows.Next() {
		var id models.UserId
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *SQLiteRoomMemberRepo) GetRoomsByUserId(ctx context.Context, userId models.UserId, limit int, cursor *string) ([]*models.Room, *string, error) {
	query := `SELECT r.id, r.name, r.created_at, r.updated_at
	FROM rooms r
	JOIN room_members rm ON r.id = rm.room_id
	WHERE rm.user_id = ?`

	args := []any{userId}

	if cursor != nil && *cursor != "" {
		query += " AND r.id < ? "
		args = append(args, *cursor)
	}

	query += " ORDER BY r.id DESC LIMIT ?"
	args = append(args, limit+1)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var rooms []*models.Room
	for rows.Next() {
		room := &models.Room{}
		err := rows.Scan(&room.Id, &room.Name, &room.CreatedAt, &room.UpdatedAt)
		if err != nil {
			return nil, nil, err
		}
		rooms = append(rooms, room)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
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
