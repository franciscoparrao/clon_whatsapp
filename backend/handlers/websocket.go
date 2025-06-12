package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   string
	username string
}

type Message struct {
	Type      string          `json:"type"`
	UserID    string          `json:"userId"`
	ChatID    string          `json:"chatId"`
	Content   string          `json:"content"`
	Timestamp int64           `json:"timestamp"`
	Data      json.RawMessage `json:"data,omitempty"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Permitir cualquier origen para desarrollo y ngrok
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client %s connected", client.userID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.mu.Unlock()
				log.Printf("Client %s disconnected", client.userID)
			} else {
				h.mu.Unlock()
			}

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastToChat sends a message to all clients in a specific chat
func (h *Hub) BroadcastToChat(chatID string, data interface{}) {
	message, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	h.broadcast <- message
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request, userID, username string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		hub:      h,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID,
		username: username,
	}

	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Process message
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		msg.UserID = c.userID
		
		// Handle different message types
		switch msg.Type {
		case "chat_message":
			// Broadcast to all clients in the chat
			c.hub.broadcast <- message
		case "typing":
			// Handle typing indicator
			c.hub.broadcast <- message
		case "read_receipt":
			// Handle read receipts
			c.hub.broadcast <- message
		default:
			log.Printf("Unknown message type: %s", msg.Type)
		}
	}
}

func (c *Client) writePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}