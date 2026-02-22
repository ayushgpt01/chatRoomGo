package user

import (
	"fmt"
	"strconv"
	"time"
)

type UserId = int64

type User struct {
	Id        UserId    `db:"id"`
	Name      string    `db:"name"`
	Username  string    `db:"user_name"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ParseUserId(id string) (UserId, error) {
	userId, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Invalid room id: %s - %s", id, err)
	}

	return userId, nil
}
