package ws

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ayushgpt01/chatRoomGo/internal/room"
	"github.com/ayushgpt01/chatRoomGo/internal/types"
)

type Hub struct {
	ctx   context.Context
	rooms map[room.RoomId]*Room
	mu    sync.RWMutex
}

func NewHub(ctx context.Context) *Hub {
	return &Hub{
		ctx:   ctx,
		rooms: make(map[room.RoomId]*Room),
		mu:    sync.RWMutex{},
	}
}

func (hub *Hub) AddRoom(roomId room.RoomId) room.RoomId {
	roomCtx, roomCancel := context.WithCancel(hub.ctx)

	room := &Room{
		id:         roomId,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan types.ChatEvent),
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

func (hub *Hub) DeleteRoom(id room.RoomId) error {
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

func (hub *Hub) GetRoom(id room.RoomId) (*Room, error) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	if hub.rooms[id] == nil {
		return nil, fmt.Errorf("No room found with id %d", id)
	}

	return hub.rooms[id], nil
}

func (hub *Hub) RoomExists(id room.RoomId) bool {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	return hub.rooms[id] != nil
}

func (hub *Hub) RegisterClient(roomId room.RoomId, client *Client) error {
	hub.mu.RLock()
	room := hub.rooms[roomId]
	hub.mu.RUnlock()

	if room == nil {
		return fmt.Errorf("No room found with id %d", roomId)
	}

	room.register <- client
	return nil
}

func (hub *Hub) UnregisterClient(roomId room.RoomId, client *Client) error {
	hub.mu.RLock()
	room := hub.rooms[roomId]
	hub.mu.RUnlock()

	if room == nil {
		return fmt.Errorf("No room found with id %d", roomId)
	}

	room.unregister <- client
	return nil
}

func (hub *Hub) Broadcast(roomId room.RoomId, evt types.ChatEvent) error {
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
