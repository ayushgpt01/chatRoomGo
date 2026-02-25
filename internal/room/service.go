package room

import (
	"context"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
)

type RoomService struct {
	roomMemberStore RoomMemberStore
	roomStore       RoomStore
	authService     *auth.AuthService
}

func NewRoomService(roomMemberStore RoomMemberStore, roomStore RoomStore, authService *auth.AuthService) *RoomService {
	return &RoomService{roomMemberStore, roomStore, authService}
}

func (srv *RoomService) HandleJoinRoom(ctx context.Context, payload JoinRoomPayload) (JoinRoomResponse, error) {
	var loginRes *auth.LoginResponse
	targetUserId := payload.UserId

	if targetUserId == 0 {
		login, err := srv.authService.HandleGuestSignup(ctx)
		if err != nil {
			return JoinRoomResponse{}, err
		}

		loginRes = &login
		targetUserId = login.User.Id
	}

	room, err := srv.roomStore.GetById(ctx, payload.Id)
	if err != nil {
		return JoinRoomResponse{}, err
	}

	if err = srv.roomMemberStore.JoinRoom(ctx, room.Id, targetUserId); err != nil {
		return JoinRoomResponse{}, err
	}

	return JoinRoomResponse{
		Room: ResponseRoom{
			Id:   room.Id,
			Name: room.Name,
		},
		Login: loginRes,
	}, nil
}

func (srv *RoomService) HandleLeaveRoom(ctx context.Context, payload LeaveRoomPayload) error {
	err := srv.roomMemberStore.LeaveRoom(ctx, payload.Id, payload.UserId)
	return err
}

func (srv *RoomService) HandleCreateRoom(ctx context.Context, payload CreateRoomPayload) (CreateRoomResponse, error) {
	return CreateRoomResponse{}, nil
}
