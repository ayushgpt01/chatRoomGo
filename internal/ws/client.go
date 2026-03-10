package ws

import (
	"encoding/json"
	"time"

	"github.com/ayushgpt01/chatRoomGo/internal/event"
	"github.com/ayushgpt01/chatRoomGo/internal/logger"
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
	log := logger.WithUserID(int(c.id)).With("room_id", int(c.roomID))

	log.Info("websocket_read_pump_started")

	c.conn.SetReadLimit(512 * 1024)
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		log.Debug("websocket_pong_received")
		return nil
	})

	defer func() {
		log.Info("websocket_read_pump_ending")
		hub.UnregisterClient(c.roomID, c)
		c.conn.Close()
	}()

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			log.Debug("websocket_read_error", "error", err.Error())
			return
		}

		log.Debug("websocket_message_received",
			"message_size", len(data),
			"payload", string(data),
		)

		var msg models.IncomingEvent

		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Warn("websocket_json_unmarshal_failed",
				"error", err.Error(),
				"payload", string(data),
			)
			continue
		}

		log.Debug("websocket_event_parsed",
			"event_type", msg.Type,
			"payload_size", len(data),
		)

		evt, err := chatService.HandleIncoming(hub.ctx, c.roomID, c.id, msg)
		if err != nil {
			log.Warn("websocket_handle_incoming_failed",
				"error", err.Error(),
				"event_type", msg.Type,
			)
			c.send <- event.NewErrorEvent(err)
			continue
		}

		log.Debug("websocket_broadcasting_event",
			"event_type", evt.Type(),
			"room_id", c.roomID,
		)

		hub.Broadcast(c.roomID, evt)
	}
}

func (c *Client) writePump() {
	log := logger.WithUserID(int(c.id)).With("room_id", int(c.roomID))

	log.Info("websocket_write_pump_started")

	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		log.Info("websocket_write_pump_ending")
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case evt, ok := <-c.send:
			if !ok {
				log.Debug("websocket_send_channel_closed")
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			outgoingEvent := models.OutgoingEvent{
				Type:    evt.Type(),
				Payload: evt.Payload(),
			}

			// Calculate payload size by marshaling to JSON
			payloadSize := 0
			if payloadBytes, err := json.Marshal(outgoingEvent.Payload); err == nil {
				payloadSize = len(payloadBytes)
			}

			log.Debug("websocket_sending_message",
				"event_type", evt.Type(),
				"payload_size", payloadSize,
			)

			if err := c.conn.WriteJSON(outgoingEvent); err != nil {
				log.Error("websocket_write_failed",
					"error", err.Error(),
					"event_type", evt.Type(),
				)
				return
			}

			log.Debug("websocket_message_sent",
				"event_type", evt.Type(),
			)

		case <-ticker.C:
			log.Debug("websocket_sending_ping")
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error("websocket_ping_failed",
					"error", err.Error(),
				)
				return
			}
		}
	}
}
