package message

import (
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/user"
)

type MessageId = int64

type Message struct {
	Id        MessageId   `db:"id"`
	Content   string      `db:"content"`
	UserId    user.UserId `db:"user_id"`
	CreatedAt time.Time   `db:"created_at"`
	UpdatedAt time.Time   `db:"updated_at"`
}
