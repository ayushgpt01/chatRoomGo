package seed

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ayushgpt01/chatRoomGo/internal/logger"
	"github.com/ayushgpt01/chatRoomGo/utils"
)

func SeedChatData(ctx context.Context, db *sql.DB) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prevent duplicate seeding
	var userCount int
	if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount); err != nil {
		logger.Error("Failed to check existing users", "error", err)
		return err
	}
	if userCount > 0 {
		logger.Info("Seed skipped: users already exist")
		return nil
	}

	passwordHash, _ := utils.HashPassword("password")

	// ---- Insert Users ----
	// Use RETURNING id and Scan for Postgres
	var aliceID, bobID int64

	err = tx.QueryRowContext(ctx, `
        INSERT INTO users (name, user_name, password_hash, account_role)
        VALUES ($1, $2, $3, $4) RETURNING id
    `, "Alice", "alice", passwordHash, "user").Scan(&aliceID)
	if err != nil {
		return fmt.Errorf("seeding Alice: %w", err)
	}

	err = tx.QueryRowContext(ctx, `
        INSERT INTO users (name, user_name, password_hash, account_role)
        VALUES ($1, $2, $3, $4) RETURNING id
    `, "Bob", "bob", passwordHash, "user").Scan(&bobID)
	if err != nil {
		return fmt.Errorf("seeding Bob: %w", err)
	}

	// ---- Insert Room ----
	var roomID int64
	err = tx.QueryRowContext(ctx, `
        INSERT INTO rooms (name)
        VALUES ($1) RETURNING id
    `, "General").Scan(&roomID)
	if err != nil {
		return fmt.Errorf("seeding room: %w", err)
	}

	// ---- Join Users to Room ----
	_, err = tx.ExecContext(ctx, `
        INSERT INTO room_members (room_id, user_id)
        VALUES ($1, $2), ($1, $3)
    `, roomID, aliceID, bobID)
	if err != nil {
		return fmt.Errorf("seeding members: %w", err)
	}

	// ---- Insert Messages ----
	_, err = tx.ExecContext(ctx, `
        INSERT INTO messages (content, user_id, room_id)
        VALUES 
            ($1, $2, $3),
            ($4, $5, $6),
            ($7, $8, $9),
            ($10, $11, $12)
    `,
		"Hey Bob 👋", aliceID, roomID,
		"Hey Alice! What's up?", bobID, roomID,
		"Just testing the seeded chat.", aliceID, roomID,
		"Looks good to me 👍", bobID, roomID,
	)
	if err != nil {
		return fmt.Errorf("seeding messages: %w", err)
	}

	logger.Info("Database seeded successfully",
		"users_created", 2,
		"rooms_created", 1,
		"messages_created", 4,
	)
	return tx.Commit()
}
