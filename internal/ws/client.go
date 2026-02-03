package ws

import (
	"github.com/ayushgpt01/chatRoomGo/internal/chat"
	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	roomID string
	conn   *websocket.Conn
	send   chan []byte
}

func (c *Client) readPump(hub *Hub, chatService *chat.ChatService) {
	defer func() {
		hub.UnregisterClient(c.roomID, c)
		c.conn.Close()
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}

		chatService.HandleIncoming(c.roomID, c.id, data)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}
