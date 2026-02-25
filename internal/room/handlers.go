package room

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
	"github.com/ayushgpt01/chatRoomGo/utils"
)

func HandleJoinRoom(srv *RoomService) http.Handler {
	type Payload struct {
		RoomId RoomId
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var temp Payload

		if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		currentUserId, _ := r.Context().Value(auth.UserIDKey).(user.UserId)

		res, err := srv.HandleJoinRoom(r.Context(), JoinRoomPayload{
			Id:     temp.RoomId,
			UserId: currentUserId,
		})

		if err != nil {
			log.Printf("POST /room/join - %v\n", err)
			if err.Error() == "invalid credentials" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, res)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}

func HandleLeaveRoom(srv *RoomService) http.Handler {
	type Payload struct {
		RoomId RoomId
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var temp Payload

		if err := json.NewDecoder(r.Body).Decode(&temp); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		currentUserId, ok := r.Context().Value(auth.UserIDKey).(user.UserId)
		if !ok {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		err := srv.HandleLeaveRoom(r.Context(), LeaveRoomPayload{
			Id:     temp.RoomId,
			UserId: currentUserId,
		})

		if err != nil {
			log.Printf("POST /room/leave - %v\n", err)
			if err.Error() == "invalid credentials" {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
