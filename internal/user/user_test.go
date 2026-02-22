package user

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/ayushgpt01/chatRoomGo/utils"
)

func setupTestRepo(t *testing.T) (*SQLiteUserRepo, *sql.DB) {
	t.Helper()

	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	repo, err := NewSqliteUserRepo(context.Background(), db)
	if err != nil {
		t.Fatalf("failed to init repo: %v", err)
	}

	return repo, db
}

func TestCreateAndGetUser(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()
	passwordHash, err := utils.HashPassword("Password")
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	userID, err := repo.Create(ctx, "raiden", "Raiden", passwordHash)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	u, err := repo.GetById(ctx, userID)
	if err != nil {
		t.Fatalf("get by id failed: %v", err)
	}

	if u.Username != "raiden" {
		t.Fatalf("expected username raiden, got %s", u.Username)
	}

	u2, err := repo.GetByUsername(ctx, "raiden")
	if err != nil {
		t.Fatalf("get by username failed: %v", err)
	}

	if u2.Id != userID {
		t.Fatalf("expected same user id")
	}
}

func TestCreateDuplicateUsernameFails(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()
	passwordHash, err := utils.HashPassword("Password")
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	_, err = repo.Create(ctx, "raiden", "Raiden", passwordHash)
	if err != nil {
		t.Fatalf("first create failed: %v", err)
	}

	passwordHash, err = utils.HashPassword("Password2")
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	_, err = repo.Create(ctx, "raiden", "Another", passwordHash)
	if err == nil {
		t.Fatalf("expected unique constraint error")
	}
}

func TestUpdateNameUpdatesTimestamp(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	passwordHash, err := utils.HashPassword("Password")
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	userID, err := repo.Create(ctx, "raiden", "Old Name", passwordHash)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	u1, err := repo.GetById(ctx, userID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	time.Sleep(1 * time.Second)

	if err := repo.UpdateName(ctx, userID, "New Name"); err != nil {
		t.Fatalf("update failed: %v", err)
	}

	u2, err := repo.GetById(ctx, userID)
	if err != nil {
		t.Fatalf("get after update failed: %v", err)
	}

	if u2.Name != "New Name" {
		t.Fatalf("name not updated")
	}

	if !u2.UpdatedAt.After(u1.UpdatedAt) {
		t.Fatalf("expected updated_at to change")
	}
}

func TestUpdateUsername(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	passwordHash, err := utils.HashPassword("Password")
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	userID, err := repo.Create(ctx, "raiden", "Raiden", passwordHash)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := repo.UpdateUsername(ctx, userID, "raiden2"); err != nil {
		t.Fatalf("update username failed: %v", err)
	}

	u, err := repo.GetById(ctx, userID)
	if err != nil {
		t.Fatalf("get failed: %v", err)
	}

	if u.Username != "raiden2" {
		t.Fatalf("username not updated")
	}
}

func TestDeleteUser(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	passwordHash, err := utils.HashPassword("Password")
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	userID, err := repo.Create(ctx, "raiden", "Raiden", passwordHash)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}

	if err := repo.DeleteById(ctx, userID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	_, err = repo.GetById(ctx, userID)
	if err == nil {
		t.Fatalf("expected error after delete")
	}
}

func TestUserNotFoundErrors(t *testing.T) {
	repo, _ := setupTestRepo(t)
	ctx := context.Background()

	_, err := repo.GetById(ctx, 999)
	if err == nil {
		t.Fatalf("expected not found error")
	}

	_, err = repo.GetByUsername(ctx, "does-not-exist")
	if err == nil {
		t.Fatalf("expected not found error")
	}

	err = repo.UpdateName(ctx, 999, "test")
	if err == nil {
		t.Fatalf("expected not found error on update")
	}

	err = repo.DeleteById(ctx, 999)
	if err == nil {
		t.Fatalf("expected not found error on delete")
	}
}
