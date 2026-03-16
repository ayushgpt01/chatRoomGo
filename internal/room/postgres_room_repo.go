package room

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type PostgresRoomRepo struct {
	db *sql.DB
}

func NewPostgresRoomRepo(ctx context.Context, db *sql.DB) *PostgresRoomRepo {
	return &PostgresRoomRepo{db}
}

func (s *PostgresRoomRepo) Create(ctx context.Context, name string) (*models.Room, error) {
	var room models.Room

	query := "INSERT INTO rooms(name) VALUES($1) RETURNING id, name, created_at, updated_at"

	err := s.db.QueryRowContext(ctx, query, name).Scan(&room.Id, &room.Name, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("inserting room with name %q: %w", name, err)
	}

	return &room, nil
}

func (s *PostgresRoomRepo) GetById(ctx context.Context, roomId models.RoomId) (*models.Room, error) {
	var room models.Room

	row := s.db.QueryRowContext(ctx, `
		SELECT r.id, r.name, r.created_at, r.updated_at
		FROM rooms r 
		WHERE r.id = $1
	`, roomId)

	err := row.Scan(&room.Id, &room.Name, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("getting room by id %d: %w", roomId, models.ErrNotFound)
		}
		return nil, fmt.Errorf("scanning room by id %d: %w", roomId, err)
	}

	return &room, nil
}

func (s *PostgresRoomRepo) UpdateName(ctx context.Context, roomId models.RoomId, name string) error {
	res, err := s.db.ExecContext(ctx, "UPDATE rooms SET name = $1 WHERE id = $2", name, roomId)
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

func (s *PostgresRoomRepo) Delete(ctx context.Context, roomId models.RoomId) error {
	res, err := s.db.ExecContext(ctx, "DELETE FROM rooms WHERE id = $1", roomId)
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
