package ws

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ayushgpt01/chatRoomGo/internal/models"
)

type Hub struct {
	ctx   context.Context
	rooms map[models.RoomId]*Room
	mu    sync.RWMutex
}

func NewHub(ctx context.Context) *Hub {
	return &Hub{
		ctx:   ctx,
		rooms: make(map[models.RoomId]*Room),
		mu:    sync.RWMutex{},
	}
}

func (hub *Hub) AddRoom(roomId models.RoomId) models.RoomId {
	roomCtx, roomCancel := context.WithCancel(hub.ctx)

	room := &Room{
		id:         roomId,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan models.ChatEvent),
		clients:    make(map[*Client]bool),
		ctx:        roomCtx,
		cancel:     roomCancel,
	}

	hub.mu.Lock()
	hub.rooms[roomId] = room
	hub.mu.Unlock()

	go room.Run()

	return roomId
}

func (hub *Hub) DeleteRoom(id models.RoomId) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	room := hub.rooms[id]
	if room == nil {
		return fmt.Errorf("No room found with id %d", id)
	}

	room.cancel()
	delete(hub.rooms, id)
	return nil
}

func (hub *Hub) GetRoom(id models.RoomId) (*Room, error) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	if hub.rooms[id] == nil {
		return nil, fmt.Errorf("No room found with id %d", id)
	}

	return hub.rooms[id], nil
}

func (hub *Hub) RoomExists(id models.RoomId) bool {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	return hub.rooms[id] != nil
}

func (hub *Hub) RegisterClient(roomId models.RoomId, client *Client) error {
	hub.mu.RLock()
	room := hub.rooms[roomId]
	hub.mu.RUnlock()

	if room == nil {
		return fmt.Errorf("No room found with id %d", roomId)
	}

	room.register <- client
	return nil
}

func (hub *Hub) UnregisterClient(roomId models.RoomId, client *Client) error {
	hub.mu.RLock()
	room := hub.rooms[roomId]
	hub.mu.RUnlock()

	if room == nil {
		return fmt.Errorf("No room found with id %d", roomId)
	}

	room.unregister <- client
	return nil
}

func (hub *Hub) Broadcast(roomId models.RoomId, evt models.ChatEvent) error {
	hub.mu.RLock()
	room := hub.rooms[roomId]
	hub.mu.RUnlock()

	if room == nil {
		return fmt.Errorf("No room found with id %d", roomId)
	}

	room.broadcast <- evt
	return nil
}

func (hub *Hub) Cleanup() {
	<-hub.ctx.Done()
	log.Printf("Cleaning up Hub...")
	hub.mu.Lock()
	defer hub.mu.Unlock()

	for id := range hub.rooms {
		delete(hub.rooms, id)
	}
}
