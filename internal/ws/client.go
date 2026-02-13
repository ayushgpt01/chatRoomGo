package ws

import (
	"encoding/json"

	"github.com/ayushgpt01/chatRoomGo/internal/chat"
	"github.com/ayushgpt01/chatRoomGo/internal/dto"
	"github.com/gorilla/websocket"
)

type Client struct {
	id     string
	roomID string
	conn   *websocket.Conn
	send   chan chat.ChatEvent
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

		var msg dto.IncomingMessage

		err = json.Unmarshal(data, &msg)
		if err != nil {
			continue
		}

		evt, err := chatService.HandleIncoming(hub.ctx, c.roomID, c.id, msg)
		if err != nil {
			c.send <- chat.NewErrorEvent(err)
			continue
		}

		hub.Broadcast(c.roomID, evt)
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for evt := range c.send {
		out := dto.OutgoingEvent{
			Type:    evt.Type(),
			Payload: evt.Payload(),
		}

		data, err := json.Marshal(out)
		if err != nil {
			continue
		}

		if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}
