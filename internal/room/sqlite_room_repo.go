package room

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
	_ "modernc.org/sqlite"
)

type RoomStore interface {
	Create(ctx context.Context, name string) (*models.Room, error)
	GetById(ctx context.Context, roomId models.RoomId) (*models.Room, error)
	UpdateName(ctx context.Context, roomId models.RoomId, name string) error
	Delete(ctx context.Context, roomId models.RoomId) error
}

type SQLiteRoomRepo struct {
	db *sql.DB
}

func NewSQLiteRoomRepo(ctx context.Context, db *sql.DB) (*SQLiteRoomRepo, error) {
	store := SQLiteRoomRepo{db}

	if err := store.init(ctx); err != nil {
		return nil, fmt.Errorf("initializing rooms table: %w", err)
	}

	return &store, nil
}

func (s *SQLiteRoomRepo) init(ctx context.Context) error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS rooms(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`

	createTriggerSQL := `CREATE TRIGGER IF NOT EXISTS update_room_timestamp
	AFTER UPDATE ON rooms
	BEGIN
		UPDATE rooms SET updated_at = CURRENT_TIMESTAMP WHERE ID = old.id;
	END;`

	if _, err := s.db.ExecContext(ctx, createTableSQL); err != nil {
		return fmt.Errorf("creating rooms table: %w", err)
	}

	if _, err := s.db.ExecContext(ctx, createTriggerSQL); err != nil {
		return fmt.Errorf("creating update_room_timestamp trigger: %w", err)
	}

	return nil
}

func (s *SQLiteRoomRepo) Create(ctx context.Context, name string) (*models.Room, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO rooms(name) VALUES(?)", name)
	if err != nil {
		return nil, fmt.Errorf("inserting room with name %q: %w", name, err)
	}

	roomId, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last insert id for room %q: %w", name, err)
	}

	room, err := s.GetById(ctx, roomId)
	if err != nil {
		return nil, fmt.Errorf("fetching created room %d: %w", roomId, err)
	}

	return room, nil
}

func (s *SQLiteRoomRepo) GetById(ctx context.Context, roomId models.RoomId) (*models.Room, error) {
	var room models.Room

	row := s.db.QueryRowContext(
		ctx,
		"SELECT id, name, created_at, updated_at FROM rooms WHERE id = ?",
		roomId,
	)

	err := row.Scan(&room.Id, &room.Name, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting room by id %d: %w", roomId, models.ErrNotFound)
		}
		return nil, fmt.Errorf("scanning room by id %d: %w", roomId, err)
	}

	return &room, nil
}

func (s *SQLiteRoomRepo) UpdateName(ctx context.Context, roomId models.RoomId, name string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE rooms SET name = ? WHERE id = ?", name, roomId)
	if err != nil {
		return fmt.Errorf("updating room name for id %d: %w", roomId, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected for update room %d: %w", roomId, err)
	}

	if count == 0 {
		return fmt.Errorf("updating room name for id %d: %w", roomId, models.ErrNotFound)
	}

	return nil
}

func (s *SQLiteRoomRepo) Delete(ctx context.Context, roomId models.RoomId) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM rooms WHERE id = ?", roomId)
	if err != nil {
		return fmt.Errorf("deleting room by id %d: %w", roomId, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("checking rows affected for delete room %d: %w", roomId, err)
	}

	if count == 0 {
		return fmt.Errorf("deleting room by id %d: %w", roomId, models.ErrNotFound)
	}

	return nil
}
