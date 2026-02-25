package ws

import (
	"context"

	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/types"
)

type Room struct {
	id         room.RoomId
	register   chan *Client
	unregister chan *Client
	broadcast  chan types.ChatEvent
	clients    map[*Client]bool

	ctx    context.Context
	cancel context.CancelFunc
}

func (r *Room) Run() {
	defer r.Cleanup()

	for {
		select {
		case <-r.ctx.Done():
			return

		case client := <-r.register:
			r.clients[client] = true

		case client := <-r.unregister:
			delete(r.clients, client)
			close(client.send)

		case msg := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					// Dropping slow clients. Need to handle later
				}
			}
		}
	}
}

func (room *Room) Cleanup() {
	for client := range room.clients {
		close(client.send)
	}

	room.clients = nil
	close(room.broadcast)
	close(room.register)
	close(room.unregister)
}
