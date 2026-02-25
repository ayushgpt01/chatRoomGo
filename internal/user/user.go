package user

import (
	"fmt"
	"strconv"
	"time"
)

type UserId = int64

type AccountRole string

const (
	AccountRoleUser  AccountRole = "user"
	AccountRoleGuest AccountRole = "guest"
	AccountRoleAdmin AccountRole = "admin"
)

func (r AccountRole) IsValid() bool {
	switch r {
	case AccountRoleUser, AccountRoleGuest, AccountRoleAdmin:
		return true
	default:
		return false
	}
}

type User struct {
	Id          UserId      `db:"id"`
	Name        string      `db:"name"`
	Username    string      `db:"user_name"`
	Password    string      `db:"password"`
	AccountRole AccountRole `db:"account_role"`
	CreatedAt   time.Time   `db:"created_at"`
	UpdatedAt   time.Time   `db:"updated_at"`
}

func ParseUserId(id string) (UserId, error) {
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid room id: %s - %s", id, err)
	}

	return userId, nil
}
