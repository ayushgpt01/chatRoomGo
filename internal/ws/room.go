package ws

import (
	"context"
	"sync"
)

type Room struct {
	id         string
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	clients    map[*Client]bool
	mu         sync.RWMutex

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
			r.mu.Lock()
			r.clients[client] = true
			r.mu.Unlock()

		case client := <-r.unregister:
			r.mu.Lock()
			delete(r.clients, client)
			close(client.send)
			r.mu.Unlock()

		case msg := <-r.broadcast:
			r.mu.RLock()

			for client := range r.clients {
				select {
				case client.send <- msg:
				default:
					// Dropping slow clients. Need to handle later
				}
			}

			r.mu.RUnlock()
		}
	}
}

func (room *Room) Cleanup() {
	room.mu.Lock()
	defer room.mu.Unlock()

	for client := range room.clients {
		close(client.send)
	}

	room.clients = nil
	close(room.broadcast)
	close(room.register)
	close(room.unregister)
}
