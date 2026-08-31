package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/tajamullone106-droid/Auryvex/internal/room"
)

type Message struct {
	Type   string      `json:"type"`
	RoomID string      `json:"room_id,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

type Client struct {
	conn   *websocket.Conn
	roomID string
	send   chan []byte
}

type Hub struct {
	mu      sync.RWMutex
	rooms   map[string]map[*Client]bool
	manager *room.Manager
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHub(manager *room.Manager) *Hub {
	return &Hub{
		rooms:   make(map[string]map[*Client]bool),
		manager: manager,
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")

	if roomID == "" {
		http.Error(w, "room is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := &Client{
		conn:   conn,
		roomID: roomID,
		send:   make(chan []byte, 32),
	}

	h.mu.Lock()

	if h.rooms[roomID] == nil {
		h.rooms[roomID] = make(map[*Client]bool)
	}

	h.rooms[roomID][client] = true

	h.mu.Unlock()

	go h.writeLoop(client)
	h.readLoop(client)
}

func (h *Hub) readLoop(client *Client) {
	defer h.remove(client)

	for {
		var msg Message

		if err := client.conn.ReadJSON(&msg); err != nil {
			return
		}

		h.broadcast(client.roomID, msg)
	}
}

func (h *Hub) writeLoop(client *Client) {
	defer client.conn.Close()

	for data := range client.send {
		if err := client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return
		}
	}
}

func (h *Hub) broadcast(roomID string, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Println("websocket marshal:", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.rooms[roomID] {
		select {
		case client.send <- data:
		default:
		}
	}
}

func (h *Hub) remove(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[client.roomID]; ok {
		delete(clients, client)

		if len(clients) == 0 {
			delete(h.rooms, client.roomID)
		}
	}

	close(client.send)
}
