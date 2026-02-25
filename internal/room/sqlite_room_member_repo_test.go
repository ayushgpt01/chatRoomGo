package room

import (
	"context"
	"database/sql"
	"testing"
)

func setupTestRoomMemberRepo(t *testing.T) (*SQLiteRoomMemberRepo, *sql.DB) {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		t.Fatalf("failed to enable fk: %v", err)
	}

	// Minimal required tables
	_, err = db.Exec(`
	CREATE TABLE users (
		id INTEGER PRIMARY KEY AUTOINCREMENT
	);
	CREATE TABLE rooms (
		id INTEGER PRIMARY KEY AUTOINCREMENT
	);
	`)
	if err != nil {
		t.Fatalf("failed creating base tables: %v", err)
	}

	repo, err := NewSQLiteRoomMemberRepo(context.Background(), db)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	return repo, db
}

func createUserAndRoom(db *sql.DB) (int64, int64) {
	res, _ := db.Exec("INSERT INTO users DEFAULT VALUES")
	userID, _ := res.LastInsertId()

	res, _ = db.Exec("INSERT INTO rooms DEFAULT VALUES")
	roomID, _ := res.LastInsertId()

	return userID, roomID
}

func TestJoinRoomAndExists(t *testing.T) {
	repo, db := setupTestRoomMemberRepo(t)
	ctx := context.Background()

	userID, roomID := createUserAndRoom(db)

	if err := repo.JoinRoom(ctx, roomID, userID); err != nil {
		t.Fatalf("join failed: %v", err)
	}

	exists, err := repo.Exists(ctx, roomID, userID)
	if err != nil {
		t.Fatalf("exists failed: %v", err)
	}

	if !exists {
		t.Fatalf("expected membership to exist")
	}
}

func TestJoinRoomDuplicateFails(t *testing.T) {
	repo, db := setupTestRoomMemberRepo(t)
	ctx := context.Background()

	userID, roomID := createUserAndRoom(db)

	if err := repo.JoinRoom(ctx, roomID, userID); err != nil {
		t.Fatalf("first join failed: %v", err)
	}

	err := repo.JoinRoom(ctx, roomID, userID)
	if err == nil {
		t.Fatalf("expected duplicate join error")
	}
}

func TestLeaveRoom(t *testing.T) {
	repo, db := setupTestRoomMemberRepo(t)
	ctx := context.Background()

	userID, roomID := createUserAndRoom(db)

	if err := repo.JoinRoom(ctx, roomID, userID); err != nil {
		t.Fatalf("join failed: %v", err)
	}

	if err := repo.LeaveRoom(ctx, roomID, userID); err != nil {
		t.Fatalf("leave failed: %v", err)
	}

	exists, _ := repo.Exists(ctx, roomID, userID)
	if exists {
		t.Fatalf("membership should not exist after leave")
	}
}

func TestLeaveNonMemberFails(t *testing.T) {
	repo, db := setupTestRoomMemberRepo(t)
	ctx := context.Background()

	userID, roomID := createUserAndRoom(db)

	err := repo.LeaveRoom(ctx, roomID, userID)
	if err == nil {
		t.Fatalf("expected error when leaving without joining")
	}
}

func TestCountByRoomId(t *testing.T) {
	repo, db := setupTestRoomMemberRepo(t)
	ctx := context.Background()

	// create 2 users
	res, _ := db.Exec("INSERT INTO users DEFAULT VALUES")
	user1, _ := res.LastInsertId()

	res, _ = db.Exec("INSERT INTO users DEFAULT VALUES")
	user2, _ := res.LastInsertId()

	res, _ = db.Exec("INSERT INTO rooms DEFAULT VALUES")
	roomID, _ := res.LastInsertId()

	repo.JoinRoom(ctx, roomID, user1)
	repo.JoinRoom(ctx, roomID, user2)

	count, err := repo.CountByRoomId(ctx, roomID)
	if err != nil {
		t.Fatalf("count failed: %v", err)
	}

	if count != 2 {
		t.Fatalf("expected count 2, got %d", count)
	}
}

func TestGetByRoomId(t *testing.T) {
	repo, db := setupTestRoomMemberRepo(t)
	ctx := context.Background()

	res, _ := db.Exec("INSERT INTO users DEFAULT VALUES")
	user1, _ := res.LastInsertId()

	res, _ = db.Exec("INSERT INTO users DEFAULT VALUES")
	user2, _ := res.LastInsertId()

	res, _ = db.Exec("INSERT INTO rooms DEFAULT VALUES")
	roomID, _ := res.LastInsertId()

	repo.JoinRoom(ctx, roomID, user1)
	repo.JoinRoom(ctx, roomID, user2)

	ids, err := repo.GetByRoomId(ctx, roomID)
	if err != nil {
		t.Fatalf("get by room failed: %v", err)
	}

	if len(ids) != 2 {
		t.Fatalf("expected 2 members, got %d", len(ids))
	}
}
