package chat

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

type Hub struct {
	mutex   sync.Mutex
	clients map[string][]*Client // lobbyID -> client
}

func NewHub() *Hub {
	return &Hub{clients: make(map[string][]*Client)}
}

func (hub *Hub) Add(lobbyID string, client *Client) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	hub.clients[lobbyID] = append(hub.clients[lobbyID], client)
}

func (hub *Hub) Broadcast(lobbyID string, msg []byte) {
	hub.mutex.Lock()
	defer hub.mutex.Unlock()

	for _, client := range hub.clients[lobbyID] {
		select {
		case client.Send <- msg:
		default:
			// drop if buffer full
		}
	}
}
