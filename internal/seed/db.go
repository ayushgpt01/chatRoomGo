package seed

import (
	"context"
	"database/sql"
	"fmt"

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
		return err
	}
	if userCount > 0 {
		fmt.Println("Seed skipped: users already exist")
		return nil
	}

	passwordHash, _ := utils.HashPassword("password")

	// ---- Insert Users ----
	res, err := tx.ExecContext(ctx, `
		INSERT INTO users (name, user_name, password_hash, account_role)
		VALUES (?, ?, ?, ?)
	`, "Alice", "alice", passwordHash, "user")
	if err != nil {
		return err
	}
	aliceID, _ := res.LastInsertId()

	res, err = tx.ExecContext(ctx, `
		INSERT INTO users (name, user_name, password_hash, account_role)
		VALUES (?, ?, ?, ?)
	`, "Bob", "bob", passwordHash, "user")
	if err != nil {
		return err
	}
	bobID, _ := res.LastInsertId()

	// ---- Insert Room ----
	res, err = tx.ExecContext(ctx, `
		INSERT INTO rooms (name)
		VALUES (?)
	`, "General")
	if err != nil {
		return err
	}
	roomID, _ := res.LastInsertId()

	// ---- Join Users to Room ----
	_, err = tx.ExecContext(ctx, `
		INSERT INTO room_members (room_id, user_id)
		VALUES (?, ?), (?, ?)
	`, roomID, aliceID, roomID, bobID)
	if err != nil {
		return err
	}

	// ---- Insert Messages ----
	_, err = tx.ExecContext(ctx, `
		INSERT INTO messages (content, user_id, room_id)
		VALUES 
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?),
			(?, ?, ?)
	`,
		"Hey Bob 👋", aliceID, roomID,
		"Hey Alice! What's up?", bobID, roomID,
		"Just testing the seeded chat.", aliceID, roomID,
		"Looks good to me 👍", bobID, roomID,
	)
	if err != nil {
		return err
	}

	fmt.Println("Seeded database")
	return tx.Commit()
}
