package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/user"
	"github.com/ayushgpt01/chatRoomGo/utils"
	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	userStore user.UserStore
	authStore AuthStore
}

// TODO - Replace this with env based secure key
const SECRET_KEY = "57668466653503e05683e68560885526888829acae9abb9c1f77df60f25aa2f0"

func NewAuthService(userStore user.UserStore, authStore AuthStore) *AuthService {
	return &AuthService{userStore, authStore}
}

func (srv *AuthService) generateAccessToken(userId user.UserId) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Minute * 15).Unix(), // 15 mins
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

func (srv *AuthService) getByAccessToken(tokenString string) (user.UserId, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(float64); ok {
			return user.UserId(sub), nil
		}
	}

	return 0, fmt.Errorf("invalid claims")
}

func (srv *AuthService) generateRefreshToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (srv *AuthService) HandleSignup(ctx context.Context, payload SignupPayload) (LoginResponse, error) {
	passwordHash, err := utils.HashPassword(payload.Password)
	if err != nil {
		return LoginResponse{}, err
	}

	userId, err := srv.userStore.Create(ctx, payload.Username, payload.Name, passwordHash)
	if err != nil {
		return LoginResponse{}, err
	}

	user, err := srv.userStore.GetById(ctx, userId)
	if err != nil {
		return LoginResponse{}, err
	}

	token, err := srv.generateAccessToken(userId)
	if err != nil {
		return LoginResponse{}, err
	}

	refreshToken := srv.generateRefreshToken()
	expiry := time.Now().Add(time.Hour * 24 * 7)

	err = srv.authStore.SaveRefreshToken(ctx, userId, refreshToken, expiry)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		User: ResponseUser{
			Id:       user.Id,
			Username: user.Username,
			Name:     user.Name,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (srv *AuthService) HandleLogin(ctx context.Context, payload LoginPayload) (LoginResponse, error) {
	user, err := srv.userStore.GetByUsername(ctx, payload.Username)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	if valid := utils.CheckPasswordHash(payload.Password, user.Password); !valid {
		return LoginResponse{}, fmt.Errorf("invalid credentials")
	}

	token, err := srv.generateAccessToken(user.Id)
	if err != nil {
		return LoginResponse{}, err
	}

	refreshToken := srv.generateRefreshToken()
	expiry := time.Now().Add(time.Hour * 24 * 7)

	err = srv.authStore.SaveRefreshToken(ctx, user.Id, refreshToken, expiry)
	if err != nil {
		return LoginResponse{}, err
	}

	return LoginResponse{
		User: ResponseUser{
			Id:       user.Id,
			Username: user.Username,
			Name:     user.Name,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (srv *AuthService) HandleRefresh(ctx context.Context, providedToken string) (string, error) {
	userId, err := srv.authStore.ValidateRefreshToken(ctx, providedToken)
	if err != nil {
		return "", fmt.Errorf("session expired, please login again")
	}

	return srv.generateAccessToken(userId)
}

func (srv *AuthService) GetCurrentUser(ctx context.Context, accessToken string) (ResponseUser, error) {
	// Get user id from accessToken
	userId, err := srv.getByAccessToken(accessToken)
	if err != nil {
		return ResponseUser{}, fmt.Errorf("invalid token")
	}

	// If valid get user and return it.
	user, err := srv.userStore.GetById(ctx, userId)
	if err != nil {
		return ResponseUser{}, fmt.Errorf("invalid token")
	}

	return ResponseUser{
		Id:       user.Id,
		Username: user.Username,
		Name:     user.Name,
	}, nil
}

func (srv *AuthService) HandleLogout(ctx context.Context, refreshToken string) error {
	return srv.authStore.DeleteRefreshToken(ctx, refreshToken)
}

func (srv *AuthService) HandleCleanup(ctx context.Context) error {
	return srv.authStore.CleanupExpiredTokens(ctx)
}
