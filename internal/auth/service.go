package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
	"github.com/ayushgpt01/chatRoomGo/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (srv *AuthService) generateAccessToken(userId models.UserId) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": time.Now().Add(time.Minute * 15).Unix(), // 15 mins
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

func (srv *AuthService) getByAccessToken(tokenString string) (models.UserId, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SECRET_KEY), nil
	})

	if err != nil || !token.Valid {
		return 0, models.ErrUnauthorized
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if sub, ok := claims["sub"].(float64); ok {
			return models.UserId(sub), nil
		}
	}

	return 0, models.ErrUnauthorized
}

func (srv *AuthService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func (srv *AuthService) HandleGuestSignup(ctx context.Context) (LoginResponse, error) {
	shortId := uuid.New().String()[:8]

	payload := SignupPayload{
		Username: "guest" + shortId,
		Password: "GUEST_PASS" + uuid.New().String(),
		Name:     "Guest User " + shortId,
	}

	userId, err := srv.userStore.Create(ctx, payload.Username, payload.Name, payload.Password, models.AccountRoleGuest)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("guest signup create user: %w", err)
	}

	user, err := srv.userStore.GetById(ctx, userId)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("guest signup get user by id=%d: %w", userId, err)
	}

	token, err := srv.generateAccessToken(userId)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("guest signup generate access token user_id=%d: %w", userId, err)
	}

	refreshToken, err := srv.generateRefreshToken()
	if err != nil {
		return LoginResponse{}, fmt.Errorf("guest signup generate refresh token: %w", err)
	}

	expiry := time.Now().Add(time.Hour * 24 * 7)

	err = srv.authStore.SaveRefreshToken(ctx, userId, refreshToken, expiry)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("guest signup save refresh token user_id=%d: %w", userId, err)
	}

	return LoginResponse{
		User: ResponseUser{
			Id:          user.Id,
			Username:    user.Username,
			Name:        user.Name,
			IsAnonymous: user.AccountRole == models.AccountRoleGuest,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (srv *AuthService) HandleSignup(ctx context.Context, payload SignupPayload) (LoginResponse, error) {
	passwordHash, err := utils.HashPassword(payload.Password)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("signup hash password: %w", err)
	}

	userId, err := srv.userStore.Create(ctx, payload.Username, payload.Name, passwordHash, models.AccountRoleUser)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("signup create user: %w", err)
	}

	user, err := srv.userStore.GetById(ctx, userId)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("signup get user by id=%d: %w", userId, err)
	}

	token, err := srv.generateAccessToken(userId)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("signup generate access token user_id=%d: %w", userId, err)
	}

	refreshToken, err := srv.generateRefreshToken()
	if err != nil {
		return LoginResponse{}, fmt.Errorf("signup generate refresh token: %w", err)
	}

	expiry := time.Now().Add(time.Hour * 24 * 7)

	err = srv.authStore.SaveRefreshToken(ctx, userId, refreshToken, expiry)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("signup save refresh token user_id=%d: %w", userId, err)
	}

	return LoginResponse{
		User: ResponseUser{
			Id:          user.Id,
			Username:    user.Username,
			Name:        user.Name,
			IsAnonymous: user.AccountRole == models.AccountRoleGuest,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (srv *AuthService) HandleLogin(ctx context.Context, payload LoginPayload) (LoginResponse, error) {
	u, err := srv.userStore.GetByUsername(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return LoginResponse{}, models.ErrUnauthorized
		}

		return LoginResponse{}, fmt.Errorf("login get user username=%s: %w", payload.Username, err)
	}

	if u.AccountRole == models.AccountRoleGuest {
		return LoginResponse{}, models.ErrUnauthorized
	}

	if valid := utils.CheckPasswordHash(payload.Password, u.Password); !valid {
		return LoginResponse{}, models.ErrUnauthorized
	}

	token, err := srv.generateAccessToken(u.Id)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("login generate access token user_id=%d: %w", u.Id, err)
	}

	refreshToken, err := srv.generateRefreshToken()
	if err != nil {
		return LoginResponse{}, fmt.Errorf("login generate refresh token: %w", err)
	}

	expiry := time.Now().Add(time.Hour * 24 * 7)

	err = srv.authStore.SaveRefreshToken(ctx, u.Id, refreshToken, expiry)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("login save refresh token user id=%d: %w", u.Id, err)
	}

	return LoginResponse{
		User: ResponseUser{
			Id:          u.Id,
			Username:    u.Username,
			Name:        u.Name,
			IsAnonymous: u.AccountRole == models.AccountRoleGuest,
		},
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (srv *AuthService) HandleRefresh(ctx context.Context, providedToken string) (string, error) {
	userId, err := srv.authStore.ValidateRefreshToken(ctx, providedToken)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return "", models.ErrUnauthorized
		}

		return "", fmt.Errorf("handle refresh validate token: %w", err)
	}

	token, err := srv.generateAccessToken(userId)
	if err != nil {
		return "", fmt.Errorf("handle refresh generate access token user_id=%d: %w", userId, err)
	}

	return token, nil
}

func (srv *AuthService) GetCurrentUser(ctx context.Context, accessToken string) (ResponseUser, error) {
	userId, err := srv.getByAccessToken(accessToken)
	if err != nil {
		if errors.Is(err, models.ErrUnauthorized) {
			return ResponseUser{}, models.ErrUnauthorized
		}

		return ResponseUser{}, fmt.Errorf("getting user id by access token: %w", err)
	}

	user, err := srv.userStore.GetById(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return ResponseUser{}, models.ErrUnauthorized
		}

		return ResponseUser{}, fmt.Errorf("getting user by id=%d: %w", userId, err)
	}

	return ResponseUser{
		Id:          user.Id,
		Username:    user.Username,
		Name:        user.Name,
		IsAnonymous: user.AccountRole == models.AccountRoleGuest,
	}, nil
}

func (srv *AuthService) HandleLogout(ctx context.Context, refreshToken string) error {
	if err := srv.authStore.DeleteRefreshToken(ctx, refreshToken); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil
		}

		return fmt.Errorf("logout delete refresh token: %w", err)
	}
	return nil
}

func (srv *AuthService) HandleCleanup(ctx context.Context) error {
	if err := srv.authStore.CleanupExpiredTokens(ctx); err != nil {
		return fmt.Errorf("cleanup expired tokens: %w", err)
	}
	return nil
}
