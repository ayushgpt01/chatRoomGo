package auth

import "github.com/ayushgpt01/chatRoomGo/internal/models"

type SignupPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	User         models.ResponseUser `json:"user"`
	Token        string              `json:"token"`
	RefreshToken string              `json:"refreshToken"`
}
