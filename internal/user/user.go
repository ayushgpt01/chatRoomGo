package user

import "time"

type UserId = int64

type User struct {
	Id        UserId    `db:"id"`
	Name      string    `db:"name"`
	Username  string    `db:"user_name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
