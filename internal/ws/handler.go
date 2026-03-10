package ws

import (
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/event"
	"github.com/ayushgpt01/chatRoomGo/internal/logger"
	"github.com/ayushgpt01/chatRoomGo/internal/models"
	"github.com/gorilla/websocket"
)

type Wshandler struct {
	hub         *Hub
	chatService *event.EventService
}

func NewWSHandler(hub *Hub, chatService *event.EventService) *Wshandler {
	return &Wshandler{
		hub,
		chatService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Wshandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomIdStr := r.URL.Query().Get("room")
	userIdStr := r.URL.Query().Get("user")
	remoteAddr := r.RemoteAddr
	userAgent := r.UserAgent()

	logger.Info("websocket_connection_attempt",
		"room_id", roomIdStr,
		"user_id", userIdStr,
		"remote_addr", remoteAddr,
		"user_agent", userAgent,
	)

	userID, err := models.ParseUserId(userIdStr)
	if err != nil {
		logger.Warn("websocket_invalid_user_id",
			"user_id_str", userIdStr,
			"error", err.Error(),
			"remote_addr", remoteAddr,
		)
		http.Error(w, "room and user required", http.StatusBadRequest)
		return
	}

	roomId, err := models.ParseRoomId(roomIdStr)
	if err != nil {
		logger.Warn("websocket_invalid_room_id",
			"room_id_str", roomIdStr,
			"error", err.Error(),
			"remote_addr", remoteAddr,
		)
		http.Error(w, "room and user required", http.StatusBadRequest)
		return
	}

	// upgrade
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error("websocket_upgrade_failed",
			"error", err.Error(),
			"user_id", userID,
			"room_id", roomId,
			"remote_addr", remoteAddr,
		)
		return
	}

	logger.Info("websocket_connection_established",
		"user_id", userID,
		"room_id", roomId,
		"remote_addr", remoteAddr,
	)

	defer ws.Close()

	// create client
	client := Client{
		id:     userID,
		conn:   ws,
		send:   make(chan models.ChatEvent),
		roomID: roomId,
	}

	// register client in room
	err = h.hub.RegisterClient(roomId, &client)
	if err != nil {
		logger.Error("websocket_register_client_failed",
			"error", err.Error(),
			"user_id", userID,
			"room_id", roomId,
		)
		ws.Close()
		return
	}

	logger.Info("websocket_client_registered",
		"user_id", userID,
		"room_id", roomId,
	)

	// start pumps
	go client.writePump()
	client.readPump(h.hub, h.chatService)
}
