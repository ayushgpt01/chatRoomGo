package room

import (
	"fmt"
	"strconv"
	"time"
)

type RoomId = int64

type Room struct {
	Id        RoomId    `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func ParseRoomId(id string) (RoomId, error) {
	roomId, err := strconv.ParseInt(id, 10, 64)
	return roomId, fmt.Errorf("Invalid room id: %s - %s", id, err)
}
