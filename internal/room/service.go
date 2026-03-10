package room

import (
	"context"
	"errors"
	"fmt"

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
			return JoinRoomResponse{}, fmt.Errorf("join room guest signup: %w", err)
		}

		loginRes = &login
		targetUserId = login.User.Id
	}

	room, err := srv.roomStore.GetById(ctx, payload.Id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return JoinRoomResponse{}, models.ErrForbidden
		}

		return JoinRoomResponse{}, fmt.Errorf("join room get room=%d: %w", payload.Id, err)
	}

	// Check if room has reached the 50 member limit
	members, err := srv.roomMemberStore.GetRoomMembers(ctx, room.Id)
	if err != nil {
		return JoinRoomResponse{}, fmt.Errorf("join room get members=%d: %w", payload.Id, err)
	}

	if len(members) >= 50 {
		return JoinRoomResponse{}, fmt.Errorf("room has reached maximum capacity of 50 members")
	}

	if err = srv.roomMemberStore.JoinRoom(ctx, room.Id, targetUserId); err != nil {
		return JoinRoomResponse{}, fmt.Errorf("join room: %w", err)
	}

	srv.hub.Broadcast(room.Id, &models.UserJoinedRoomEvent{
		Data: models.UserJoinedRoomPayload{
			RoomID: room.Id,
			UserID: targetUserId,
		},
	})

	return JoinRoomResponse{
		Room: ResponseRoom{
			Id:               room.Id,
			Name:             room.Name,
			ParticipantCount: 0,
			UpdatedAt:        room.UpdatedAt,
			Members:          []*models.ResponseUser{},
		},
		Login: loginRes,
	}, nil
}

func (srv *RoomService) HandleLeaveRoom(ctx context.Context, payload LeaveRoomPayload) error {
	err := srv.roomMemberStore.LeaveRoom(ctx, payload.Id, payload.UserId)

	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.ErrForbidden
		}

		return fmt.Errorf("handle leave room: %w", err)
	}

	srv.hub.Broadcast(payload.Id, &models.UserLeftRoomEvent{
		Data: models.UserLeftRoomPayload{
			RoomID: payload.Id,
			UserID: payload.UserId,
		},
	})

	return nil
}

func (srv *RoomService) HandleCreateRoom(ctx context.Context, payload CreateRoomPayload) (CreateRoomResponse, error) {
	room, err := srv.roomStore.Create(ctx, payload.Name)
	if err != nil {
		return CreateRoomResponse{}, fmt.Errorf("create room name=%s: %w", payload.Name, err)
	}

	if err = srv.roomMemberStore.JoinRoom(ctx, room.Id, payload.UserId); err != nil {
		return CreateRoomResponse{}, fmt.Errorf("create room join user: %w", err)
	}

	return CreateRoomResponse{
		Room: ResponseRoom{
			Id:               room.Id,
			Name:             room.Name,
			ParticipantCount: 0,
			UpdatedAt:        room.UpdatedAt,
			Members:          []*models.ResponseUser{},
		},
	}, nil
}

func (srv *RoomService) HandleGetRooms(ctx context.Context, payload GetRoomPayload) (GetRoomResponse, error) {
	rooms, nextCursor, err := srv.roomMemberStore.GetRoomsByUserId(ctx, payload.UserId, payload.Limit, payload.Cursor)
	if err != nil {
		return GetRoomResponse{}, fmt.Errorf("get rooms: %w", err)
	}

	var responseRooms []ResponseRoom
	for _, room := range rooms {
		members, err := srv.roomMemberStore.GetRoomMembers(ctx, room.Id)
		if err != nil {
			// Don't fail the entire operation if we can't get members
			members = []*models.User{}
		}

		responseMembers := make([]*models.ResponseUser, len(members))
		for i, member := range members {
			responseMembers[i] = &models.ResponseUser{
				Id:          member.Id,
				Username:    member.Username,
				Name:        member.Name,
				IsAnonymous: member.AccountRole == models.AccountRoleGuest,
			}
		}

		responseRooms = append(responseRooms, ResponseRoom{
			Id:               room.Id,
			Name:             room.Name,
			ParticipantCount: len(responseMembers),
			UpdatedAt:        room.UpdatedAt,
			Members:          responseMembers,
		})
	}

	return GetRoomResponse{
		Rooms:      responseRooms,
		NextCursor: nextCursor,
	}, nil
}
