package ws

import (
	"encoding/json"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/event"
	"github.com/ayushgpt01/chatRoomGo/internal/models"
	"github.com/gorilla/websocket"
)

type Client struct {
	id     models.UserId
	roomID models.RoomId
	conn   *websocket.Conn
	send   chan models.ChatEvent
}

type connectedEvent struct {
}

func (e *connectedEvent) Type() string {
	return "connected"
}

func (e *connectedEvent) Payload() any {
	return nil
}

func (c *Client) readPump(hub *Hub, chatService *event.EventService) {
	c.conn.SetReadLimit(512 * 1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	c.send <- &connectedEvent{}

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
			c.send <- event.NewErrorEvent(err)
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

			c.conn.WriteJSON(models.OutgoingEvent{
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
