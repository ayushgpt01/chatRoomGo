package room

import (
	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
)

type JoinRoomPayload struct {
	Id     RoomId      `json:"roomId"`
	UserId user.UserId `json:"userId"`
}

type ResponseRoom struct {
	Id   RoomId `json:"id"`
	Name string `json:"name"`
}

type JoinRoomResponse struct {
	Room  ResponseRoom        `json:"room"`
	Login *auth.LoginResponse `json:"login,omitempty"`
}

type LeaveRoomPayload struct {
	Id     RoomId      `json:"roomId"`
	UserId user.UserId `json:"userId"`
}

type CreateRoomPayload struct {
	UserId user.UserId `json:"userId"`
}

type CreateRoomResponse struct {
	Room ResponseRoom `json:"room"`
}
