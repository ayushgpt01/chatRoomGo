package message

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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

type request struct {
	Content string  `json:"content"`
	Nonce   *string `json:"nonce,omitempty"`
}

func (s request) Valid(ctx context.Context) map[string]string {
	problems := map[string]string{}
	if strings.TrimSpace(s.Content) == "" {
		problems["content"] = "Content cannot be empty"
	}
	if len(s.Content) > 2000 {
		problems["content"] = "Content too long"
	}
	return problems
}

func HandleSendMessage(srv *MessageService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		roomId, err := models.ParseRoomId(r.PathValue("roomId"))
		if err != nil {
			http.Error(w, "Invalid room id", http.StatusBadRequest)
			return
		}

		currentUserId, ok := r.Context().Value(auth.UserIDKey).(models.UserId)
		if !ok {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		body, ok := utils.HandleDecode[request](w, r)
		if !ok {
			return
		}

		res, err := srv.HandleSendMessage(r.Context(), SendMessagePayload{
			UserId:  currentUserId,
			RoomId:  roomId,
			Content: body.Content,
			Nonce:   body.Nonce,
		})

		if err != nil {
			utils.HandleServiceError(w, "POST /room/{roomId}/messages", err)
			return
		}

		if err := utils.Encode(w, r, http.StatusCreated, res); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}

func HandleEditMessage(srv *MessageService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roomId, err := models.ParseRoomId(r.PathValue("roomId"))
		if err != nil {
			http.Error(w, "Invalid room id", http.StatusBadRequest)
			return
		}

		messageId, err := models.ParseMessageId(r.PathValue("messageId"))
		if err != nil {
			http.Error(w, "Invalid message id", http.StatusBadRequest)
			return
		}

		currentUserId, ok := r.Context().Value(auth.UserIDKey).(models.UserId)
		if !ok {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		body, ok := utils.HandleDecode[request](w, r)
		if !ok {
			return
		}

		res, err := srv.HandleEditMessage(r.Context(), EditMessagePayload{
			UserId:    currentUserId,
			MessageId: messageId,
			RoomId:    roomId,
			Content:   body.Content,
		})

		if err != nil {
			utils.HandleServiceError(w, "PATCH /rooms/{roomId}/messages/{id}", err)
			return
		}

		if err := utils.Encode(w, r, http.StatusOK, res); err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
		}
	})
}

func HandleDeleteMessage(srv *MessageService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		roomId, err := models.ParseRoomId(r.PathValue("roomId"))
		if err != nil {
			http.Error(w, "Invalid room id", http.StatusBadRequest)
			return
		}

		messageId, err := models.ParseMessageId(r.PathValue("messageId"))
		if err != nil {
			http.Error(w, "Invalid message id", http.StatusBadRequest)
			return
		}

		currentUserId, ok := r.Context().Value(auth.UserIDKey).(models.UserId)
		if !ok {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		err = srv.HandleDeleteMessage(r.Context(), DeleteMessagePayload{
			UserId:    currentUserId,
			MessageId: messageId,
			RoomId:    roomId,
		})

		if err != nil {
			utils.HandleServiceError(w, "DELETE /messages/{id}", err)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
