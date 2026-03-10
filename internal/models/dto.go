package models

import "time"

type ResponseMessage struct {
	Id         MessageId  `json:"id"`
	Content    string     `json:"content"`
	SenderName string     `json:"senderName"`
	SenderId   UserId     `json:"senderId"`
	SentAt     time.Time  `json:"sentAt"`
	EditedAt   *time.Time `json:"editedAt"`
	Nonce      *string    `json:"nonce"`
	RoomId     RoomId     `json:"roomId"`
	Delivered  bool       `json:"delivered"`
	ReadBy     []UserId   `json:"readBy"`
}

type ResponseUser struct {
	Id          UserId `json:"id"`
	Username    string `json:"username"`
	Name        string `json:"name"`
	IsAnonymous bool   `json:"isAnonymous"`
}
