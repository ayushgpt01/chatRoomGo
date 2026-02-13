package ws

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ayushgpt01/chatRoomGo/internal/chat"
	"github.com/google/uuid"
)

type Hub struct {
	ctx   context.Context
	rooms map[string]*Room
	mu    sync.RWMutex
}

func NewHub(ctx context.Context) *Hub {
	return &Hub{
		ctx:   ctx,
		rooms: make(map[string]*Room),
		mu:    sync.RWMutex{},
	}
}

func (hub *Hub) AddRoom() string {
	id := uuid.NewString()
	roomCtx, roomCancel := context.WithCancel(hub.ctx)

	room := &Room{
		id:         id,
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan chat.ChatEvent),
		clients:    make(map[*Client]bool),
		mu:         sync.RWMutex{},
		ctx:        roomCtx,
		cancel:     roomCancel,
	}

	hub.mu.Lock()
	hub.rooms[id] = room
	hub.mu.Unlock()

	go room.Run()

	return id
}

func (hub *Hub) DeleteRoom(id string) error {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	room := hub.rooms[id]
	if room == nil {
		return fmt.Errorf("No room found with id %s", id)
	}

	room.cancel()
	delete(hub.rooms, id)
	return nil
}

func (hub *Hub) GetRoom(id string) (*Room, error) {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	if hub.rooms[id] == nil {
		return nil, fmt.Errorf("No room found with id %s", id)
	}

	return hub.rooms[id], nil
}

func (hub *Hub) RoomExists(id string) bool {
	hub.mu.RLock()
	defer hub.mu.RUnlock()

	return hub.rooms[id] != nil
}

func (hub *Hub) RegisterClient(roomId string, client *Client) error {
	hub.mu.RLock()
	room := hub.rooms[roomId]
	hub.mu.RUnlock()

	if room == nil {
		return fmt.Errorf("No room found with id %s", roomId)
	}

	room.register <- client
	return nil
}

func (hub *Hub) UnregisterClient(roomId string, client *Client) error {
	hub.mu.RLock()
	room := hub.rooms[roomId]
	hub.mu.RUnlock()

	if room == nil {
		return fmt.Errorf("No room found with id %s", roomId)
	}

	room.unregister <- client
	return nil
}

func (hub *Hub) Broadcast(roomId string, evt chat.ChatEvent) error {
	hub.mu.RLock()
	room := hub.rooms[roomId]
	hub.mu.RUnlock()

	if room == nil {
		return fmt.Errorf("No room found with id %s", roomId)
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
