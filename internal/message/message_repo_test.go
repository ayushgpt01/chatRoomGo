package message

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func setupTestRepo(t *testing.T) (*SQLiteMessageRepo, *sql.DB) {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, _ = db.Exec("PRAGMA foreign_keys = ON;")

	_, err = db.Exec(`
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT
	);
	CREATE TABLE rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT
	);`)
	if err != nil {
		t.Fatalf("failed creating tables: %v", err)
	}

	repo, err := NewSQLiteMessageRepo(context.Background(), db)
	if err != nil {
		t.Fatalf("failed creating repo: %v", err)
	}

	return repo, db
}

func TestCreateAndGetMessage(t *testing.T) {
	repo, db := setupTestRepo(t)
	ctx := context.Background()

	res, _ := db.Exec("INSERT INTO users DEFAULT VALUES")
	userID, _ := res.LastInsertId()

	res, _ = db.Exec("INSERT INTO rooms DEFAULT VALUES")
	roomID, _ := res.LastInsertId()

	msgID, err := repo.Create(ctx, roomID, userID, "hello")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	msg, err := repo.GetById(ctx, msgID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if msg.Content != "hello" {
		t.Fatalf("expected content 'hello', got %s", msg.Content)
	}

	if msg.UserId != userID {
		t.Fatalf("wrong user id")
	}
}

func TestDeleteMessage(t *testing.T) {
	repo, db := setupTestRepo(t)
	ctx := context.Background()

	res, _ := db.Exec("INSERT INTO users DEFAULT VALUES")
	userID, _ := res.LastInsertId()

	res, _ = db.Exec("INSERT INTO rooms DEFAULT VALUES")
	roomID, _ := res.LastInsertId()

	msgID, _ := repo.Create(ctx, roomID, userID, "delete me")

	if err := repo.DeleteById(ctx, msgID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err := repo.GetById(ctx, msgID)
	if err == nil {
		t.Fatalf("expected error after delete")
	}
}

func TestUpdateContentUpdatesTimestamp(t *testing.T) {
	repo, db := setupTestRepo(t)
	ctx := context.Background()

	res, _ := db.Exec("INSERT INTO users DEFAULT VALUES")
	userID, _ := res.LastInsertId()

	res, _ = db.Exec("INSERT INTO rooms DEFAULT VALUES")
	roomID, _ := res.LastInsertId()

	msgID, _ := repo.Create(ctx, roomID, userID, "old")

	msg1, _ := repo.GetById(ctx, msgID)

	time.Sleep(1 * time.Second)

	if err := repo.UpdateContent(ctx, msgID, "new"); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	msg2, _ := repo.GetById(ctx, msgID)

	if !msg2.UpdatedAt.After(msg1.UpdatedAt) {
		t.Fatalf("updated_at did not change")
	}
}

func TestCreateFailsWithInvalidForeignKeys(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.Create(ctx, 999, 999, "invalid")

	if err == nil {
		t.Fatalf("expected foreign key error")
	}
}
