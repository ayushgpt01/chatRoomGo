package models

import "time"

type ResponseMessage struct {
	Id         MessageId  `json:"id"`
	Content    string     `json:"content"`
	SenderName string     `json:"senderName"`
	SenderId   UserId     `json:"senderId"`
	SentAt     time.Time  `json:"sentAt"`
	EditedAt   *time.Time `json:"editedAt"`
	Read       bool       `json:"read"`
	Nonce      *string    `json:"nonce,omitempty"`
	RoomId     RoomId     `json:"roomId"`
}
