package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// MessageLogger is an interface for logging messages
type MessageLogger interface {
	SaveMessage(content string, clientID string) error
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients       map[*Client]bool
	broadcast     chan []byte
	register      chan *Client
	unregister    chan *Client
	mu            sync.RWMutex
	messageLogger MessageLogger
}

// NewHub creates a new Hub instance
func NewHub(messageLogger MessageLogger) *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		broadcast:     make(chan []byte),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		messageLogger: messageLogger,
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			log.Printf("Client registered. Total clients: %d", len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			log.Printf("Client unregistered. Total clients: %d", len(h.clients))

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

// Register registers a new client with the hub
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister unregisters a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// BroadcastJSON broadcasts a JSON object to all clients
func (h *Hub) BroadcastJSON(data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal broadcast data: %v", err)
		return
	}
	h.Broadcast(jsonData)
}

// Client represents a WebSocket connection
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	mu   sync.Mutex
	ID   string
}

// Send sends a message to this client
func (c *Client) Send(message []byte) {
	select {
	case c.send <- message:
	default:
		close(c.send)
	}
}

// NewClient creates a new Client instance
func NewClient(hub *Hub, conn *websocket.Conn, clientID string) *Client {
	return &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		ID:   clientID,
	}
}

// ReadPump pumps messages from the WebSocket connection to the hub
func (c *Client) ReadPump() {
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

		// Save message with client ID before broadcasting
		if c.hub.messageLogger != nil {
			if err := c.hub.messageLogger.SaveMessage(string(message), c.ID); err != nil {
				log.Printf("Failed to save message from client %s: %v", c.ID, err)
			}
		}

		// Broadcast message to all clients
		c.hub.broadcast <- message
	}
}

// WritePump pumps messages from the hub to the WebSocket connection
func (c *Client) WritePump() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			c.mu.Lock()
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.mu.Unlock()
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				c.mu.Unlock()
				return
			}
			w.Write(message)

			// Add queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				c.mu.Unlock()
				return
			}
			c.mu.Unlock()
		}
	}
}
