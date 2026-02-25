package room

import (
	"context"
	"database/sql"
	"testing"
	"time"
)

func setupTestRepo(t *testing.T) (*SQLiteRoomRepo, *sql.DB) {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	_, _ = db.Exec("PRAGMA foreign_keys = ON;")

	repo, err := NewSQLiteRoomRepo(context.Background(), db)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	return repo, db
}

func TestCreateAndGetRoom(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	roomID, err := repo.Create(ctx, "general")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	r, err := repo.GetById(ctx, roomID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if r.Name != "general" {
		t.Fatalf("expected name 'general', got %s", r.Name)
	}

	if r.Id != roomID {
		t.Fatalf("expected id %d, got %d", roomID, r.Id)
	}
}

func TestUpdateNameUpdatesTimestamp(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	roomID, err := repo.Create(ctx, "old")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	r1, err := repo.GetById(ctx, roomID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	time.Sleep(1 * time.Second)

	if err := repo.UpdateName(ctx, roomID, "new"); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	r2, err := repo.GetById(ctx, roomID)
	if err != nil {
		t.Fatalf("get after update failed: %v", err)
	}

	if r2.Name != "new" {
		t.Fatalf("expected updated name 'new', got %s", r2.Name)
	}

	if !r2.UpdatedAt.After(r1.UpdatedAt) {
		t.Fatalf("expected updated_at to change")
	}
}

func TestDeleteRoom(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	roomID, err := repo.Create(ctx, "temp")
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := repo.Delete(ctx, roomID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = repo.GetById(ctx, roomID)
	if err == nil {
		t.Fatalf("expected error after delete")
	}
}

func TestRoomNotFoundErrors(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetById(ctx, 999)
	if err == nil {
		t.Fatalf("expected not found error")
	}

	err = repo.UpdateName(ctx, 999, "doesnt-exist")
	if err == nil {
		t.Fatalf("expected not found error on update")
	}

	err = repo.Delete(ctx, 999)
	if err == nil {
		t.Fatalf("expected not found error on delete")
	}
}
