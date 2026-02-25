package ws

import (
	"log"
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/auth"
	"github.com/ayushgpt01/chatRoomGo/internal/chat"
	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/types"
	"github.com/ayushgpt01/chatRoomGo/internal/user"
	"github.com/gorilla/websocket"
)

type Wshandler struct {
	hub         *Hub
	chatService *chat.ChatService
}

func NewWSHandler(hub *Hub, chatService *chat.ChatService) *Wshandler {
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
	userID, ok := r.Context().Value(auth.UserIDKey).(user.UserId)
	if !ok || userID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	roomId, err := room.ParseRoomId(roomIdStr)
	if err != nil {
		http.Error(w, "room and user required", http.StatusBadRequest)
		return
	}

	// upgrade
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	defer ws.Close()

	// create client
	client := Client{
		id:     userID,
		conn:   ws,
		send:   make(chan types.ChatEvent),
		roomID: roomId,
	}

	// register client in room
	err = h.hub.RegisterClient(roomId, &client)
	if err != nil {
		log.Println("register client error:", err)
		return
	}

	// start pumps
	go client.writePump()
	client.readPump(h.hub, h.chatService)
}
