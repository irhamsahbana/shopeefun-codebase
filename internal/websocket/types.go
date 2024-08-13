package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

// Client represents a connected client
type Client struct {
	id   string
	name string
	role string
	conn *websocket.Conn
	send chan []byte
	room *Hub
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	roomId     string
	done       chan struct{}
	clients    map[*Client]struct{}
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

// NewHub creates a new Hub
func NewHub(roomId string) *Hub {
	return &Hub{
		roomId:     roomId,
		done:       make(chan struct{}),
		clients:    make(map[*Client]struct{}),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub and listens for incoming messages
func (h *Hub) Run(rooms map[string]*Hub) {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = struct{}{}
			// Reset cleanup if a new client joins
			close(h.done)
			h.done = make(chan struct{})
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			// send message to other clients that the user has left
			message := []byte(`{
					"event": "chat-group-user-left",
					"data": {
						"user_id": "` + client.id + `",
						"name": "` + client.name + `",
						"message": "left the chat"
				}
			}`)

			if _, ok := h.clients[client]; ok {
				for c := range h.clients {
					if c != client {
						c.send <- message
					}
				}
				delete(h.clients, client)
				close(client.send)
				if len(h.clients) == 0 {
					go h.cleanup(rooms)
				}
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.Lock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// cleanup removes the hub if no clients are connected after a timeout
func (h *Hub) cleanup(rooms map[string]*Hub) {
	select {
	case <-time.After(30 * time.Second): // Timeout period
		h.mu.Lock()
		if len(h.clients) == 0 {
			delete(rooms, h.roomId)
			close(h.done)
			log.Info().Msgf("Room %s deleted due to inactivity", h.roomId)
		}
		h.mu.Unlock()
	case <-h.done:
		return
	}
}
