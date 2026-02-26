package room

import (
	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type JoinRoomPayload struct {
	Id     models.RoomId `json:"roomId"`
	UserId models.UserId `json:"userId"`
}

type ResponseRoom struct {
	Id   models.RoomId `json:"id"`
	Name string        `json:"name"`
}

type JoinRoomResponse struct {
	Room  ResponseRoom        `json:"room"`
	Login *auth.LoginResponse `json:"login,omitempty"`
}

type LeaveRoomPayload struct {
	Id     models.RoomId `json:"roomId"`
	UserId models.UserId `json:"userId"`
}

type CreateRoomPayload struct {
	UserId models.UserId `json:"userId"`
}

type CreateRoomResponse struct {
	Room ResponseRoom `json:"room"`
}
