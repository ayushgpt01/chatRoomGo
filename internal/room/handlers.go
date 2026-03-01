package room

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/models"
	"github.com/ayushgpt01/chatRoomGo/utils"
)

type Request struct {
	RoomId models.RoomId `json:"roomId"`
}

func (p Request) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)
	if p.RoomId == 0 {
		problems["roomId"] = "Room Id is required"
	}
	return problems
}

func HandleJoinRoom(srv *RoomService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp, ok := utils.HandleDecode[Request](w, r)
		if !ok {
			return
		}

		currentUserId, _ := r.Context().Value(auth.UserIDKey).(models.UserId)

		res, err := srv.HandleJoinRoom(r.Context(), JoinRoomPayload{
			Id:     temp.RoomId,
			UserId: currentUserId,
		})

		if err != nil {
			utils.HandleServiceError(w, "POST /room/join", err)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, res)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}

func HandleLeaveRoom(srv *RoomService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp, ok := utils.HandleDecode[Request](w, r)
		if !ok {
			return
		}

		currentUserId, ok := r.Context().Value(auth.UserIDKey).(models.UserId)
		if !ok {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		err := srv.HandleLeaveRoom(r.Context(), LeaveRoomPayload{
			Id:     temp.RoomId,
			UserId: currentUserId,
		})

		if err != nil {
			utils.HandleServiceError(w, "POST /room/leave", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

type CreateRoomRequest struct {
	Name string `json:"name"`
}

func (p CreateRoomRequest) Valid(ctx context.Context) map[string]string {
	problems := make(map[string]string)
	if len(p.Name) == 0 {
		problems["name"] = "room name is required"
	}
	return problems
}

func HandleCreateRoom(srv *RoomService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		temp, ok := utils.HandleDecode[CreateRoomRequest](w, r)
		if !ok {
			return
		}

		currentUserId, ok := r.Context().Value(auth.UserIDKey).(models.UserId)
		if !ok {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		res, err := srv.HandleCreateRoom(r.Context(), CreateRoomPayload{
			UserId: currentUserId,
			Name:   temp.Name,
		})

		if err != nil {
			utils.HandleServiceError(w, "POST /room/create", err)
			return
		}

		err = utils.Encode(w, r, http.StatusCreated, res)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}

func HandleGetRooms(srv *RoomService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()

		limitStr := query.Get("limit")
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			limit = 10
		}

		if limit > 100 {
			http.Error(w, "Limit should be under 100", http.StatusBadRequest)
			return
		}

		var cursor *string
		if c := query.Get("cursor"); c != "" {
			cursor = &c
		}

		currentUserId, ok := r.Context().Value(auth.UserIDKey).(models.UserId)
		if !ok {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		res, err := srv.HandleGetRooms(r.Context(), GetRoomPayload{
			UserId: currentUserId,
			Limit:  limit,
			Cursor: cursor,
		})

		if err != nil {
			utils.HandleServiceError(w, "GET /room/getAll", err)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, res)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}
