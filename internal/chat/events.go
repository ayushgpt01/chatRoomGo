package chat

import (
	"github.com/ayushgpt01/chatRoomGo/internal/message"
)

type ChatEvent interface {
	Type() string
	Payload() any
}

type NewMessageEvent struct {
	Message *message.Message
}

func (e *NewMessageEvent) Type() string {
	return "new_message"
}

func (e *NewMessageEvent) Payload() any {
	return e.Message
}
