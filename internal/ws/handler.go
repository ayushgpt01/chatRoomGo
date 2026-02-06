package ws

import (
	"log"
	"net/http"

	"github.com/ayushgpt01/chatRoomGo/internal/chat"
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
	roomId := r.URL.Query().Get("room")
	user := r.URL.Query().Get("user")

	if roomId == "" || user == "" {
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
		id:     user,
		conn:   ws,
		send:   make(chan chat.ChatEvent),
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
