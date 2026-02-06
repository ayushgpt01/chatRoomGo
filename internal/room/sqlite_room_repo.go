package room

import (
	"context"
	"database/sql"
	"fmt"
)

type SQLiteRoomRepo struct {
	db *sql.DB
}

func NewSQLiteRoomRepo(ctx context.Context, db *sql.DB) (*SQLiteRoomRepo, error) {
	store := SQLiteRoomRepo{db}

	if err := store.init(ctx); err != nil {
		return nil, err
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
		return err
	}

	if _, err := s.db.ExecContext(ctx, createTriggerSQL); err != nil {
		return err
	}

	return nil
}

func (s *SQLiteRoomRepo) Create(ctx context.Context, name string) (RoomId, error) {
	res, err := s.db.ExecContext(ctx, "INSERT INTO rooms(name) VALUES(?)", name)
	if err != nil {
		return 0, err
	}

	roomId, err := res.LastInsertId()
	return roomId, err
}

func (s *SQLiteRoomRepo) GetById(ctx context.Context, roomId RoomId) (*Room, error) {
	var room Room
	row := s.db.QueryRowContext(ctx, "SELECT id, name, created_at, updated_at FROM rooms WHERE id = ?", roomId)
	err := row.Scan(&room.Id, &room.Name, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("GetById %d: no such room", roomId)
		}

		return nil, fmt.Errorf("GetById %d: %v", roomId, err)
	}

	return &room, nil
}

func (s *SQLiteRoomRepo) UpdateName(ctx context.Context, roomId RoomId, name string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE rooms SET name = ? WHERE id = ?", name, roomId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Room not found for id: %d", roomId)
	}

	return nil
}

func (s *SQLiteRoomRepo) Delete(ctx context.Context, roomId RoomId) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM rooms WHERE id = ?", roomId)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("Room not found for id: %d", roomId)
	}

	return nil
}
