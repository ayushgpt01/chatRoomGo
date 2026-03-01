package message

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/models"
	"github.com/ayushgpt01/chatRoomGo/utils"
)

func HandleGetMessages(srv *MessageService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roomId, err := models.ParseRoomId(r.PathValue("roomId"))
		if err != nil {
			http.Error(w, "Invalid room id", http.StatusBadRequest)
			return
		}

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

		res, err := srv.HandleGetMessages(r.Context(), GetMessagesPayload{
			UserId: currentUserId,
			RoomId: roomId,
			Limit:  limit,
			Cursor: cursor,
		})

		if err != nil {
			utils.HandleServiceError(w, fmt.Sprintf("GET /room/%d/messages", roomId), err)
			return
		}

		err = utils.Encode(w, r, http.StatusOK, res)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}
