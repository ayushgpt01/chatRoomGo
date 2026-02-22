package auth

import "github.com/ayushgpt01/chatRoomGo/internal/user"

type SignupPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseUser struct {
	Id       user.UserId `json:"id"`
	Username string      `json:"username"`
	Name     string      `json:"name"`
}

type LoginResponse struct {
	User         ResponseUser `json:"user"`
	Token        string       `json:"token"`
	RefreshToken string       `json:"refreshToken"`
}
