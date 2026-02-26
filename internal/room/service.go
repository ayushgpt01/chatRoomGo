package room

import (
	"context"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type RoomService struct {
	roomMemberStore RoomMemberStore
	roomStore       RoomStore
	authService     *auth.AuthService
	hub             models.HubBroadcaster
}

func NewRoomService(roomMemberStore RoomMemberStore, roomStore RoomStore, authService *auth.AuthService, hub models.HubBroadcaster) *RoomService {
	return &RoomService{roomMemberStore, roomStore, authService, hub}
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

	srv.hub.Broadcast(room.Id, &models.BaseEvent{
		EventType: models.EventUserJoinedRoom,
		Data: map[string]any{
			"roomId": room.Id,
			"userId": targetUserId,
		},
	})

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

	if err == nil {
		srv.hub.Broadcast(payload.Id, &models.BaseEvent{
			EventType: models.EventUserLeftRoom,
			Data: map[string]any{
				"roomId": payload.Id,
				"userId": payload.UserId,
			},
		})
	}

	return err
}

func (srv *RoomService) HandleCreateRoom(ctx context.Context, payload CreateRoomPayload) (CreateRoomResponse, error) {
	return CreateRoomResponse{}, nil
}
