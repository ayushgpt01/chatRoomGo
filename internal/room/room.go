package room

import "time"

type RoomId = int64

type Room struct {
	Id        RoomId    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
