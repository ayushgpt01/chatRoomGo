package ws

import (
	"encoding/json"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/chat"
	"github.com/ayushgpt01/chatRoomGo/internal/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	id     models.UserId
	roomID models.RoomId
	conn   *websocket.Conn
	send   chan models.ChatEvent
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

		var msg models.IncomingEvent

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
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case evt, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			type OutgoingEvent struct {
				Type    string `json:"type"`
				Payload any    `json:"payload"`
			}

			c.conn.WriteJSON(OutgoingEvent{
				Type:    evt.Type(),
				Payload: evt.Payload(),
			})
		case <-ticker.C:
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
