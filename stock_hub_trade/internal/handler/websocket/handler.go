package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	"stock_hub_trade/internal/service"
	"stock_hub_trade/pkg/model"
	ws "stock_hub_trade/pkg/websocket"

	gorillaWS "github.com/gorilla/websocket"
)

var upgrader = gorillaWS.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// Handler handles WebSocket connections
type Handler struct {
	messageService *service.MessageService
	userService    *service.UserService
	orderService   *service.StockOrderService
	hub            *ws.Hub
	clientUsers    map[*ws.Client]*model.User
}

// NewHandler creates a new WebSocket handler
func NewHandler(
	messageService *service.MessageService,
	userService *service.UserService,
	orderService *service.StockOrderService,
	hub *ws.Hub,
) *Handler {
	return &Handler{
		messageService: messageService,
		userService:    userService,
		orderService:   orderService,
		hub:            hub,
		clientUsers:    make(map[*ws.Client]*model.User),
	}
}

// HandleConnection handles WebSocket connection upgrade and messaging
func (h *Handler) HandleConnection(w http.ResponseWriter, r *http.Request) {
	log.Printf("WebSocket connection attempt from %s", r.RemoteAddr)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}

	// Create client and register with hub
	client := ws.NewClient(h.hub, conn, "")
	h.hub.Register(client)

	// Start write pump
	go client.WritePump()

	// Handle messages
	go h.readPump(client, conn)
}

func (h *Handler) readPump(client *ws.Client, conn *gorillaWS.Conn) {
	defer func() {
		delete(h.clientUsers, client)
		h.hub.Unregister(client)
		conn.Close()
	}()

	for {
		_, messageBytes, err := conn.ReadMessage()
		if err != nil {
			if gorillaWS.IsUnexpectedCloseError(err, gorillaWS.CloseGoingAway, gorillaWS.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(messageBytes, &msg); err != nil {
			log.Printf("Failed to parse message: %v", err)
			continue
		}

		h.handleMessage(client, &msg)
	}
}

func (h *Handler) handleMessage(client *ws.Client, msg *WSMessage) {
	switch msg.Type {
	case "register":
		h.handleRegister(client, msg)
	case "message":
		h.handleChatMessage(client, msg)
	case "order":
		h.handleOrder(client, msg)
	case "cancel_order":
		h.handleCancelOrder(client, msg)
	}
}

func (h *Handler) handleRegister(client *ws.Client, msg *WSMessage) {
	var data RegisterData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		h.sendError(client, "Invalid registration data")
		return
	}

	// Validate required fields
	if data.Username == "" {
		h.sendError(client, "Username is required")
		return
	}
	if data.Email == "" {
		h.sendError(client, "Email is required")
		return
	}

	user, err := h.userService.RegisterUser(data.Username, data.Email)
	if err != nil {
		h.sendError(client, "Failed to register user: "+err.Error())
		return
	}

	h.clientUsers[client] = user
	client.ID = data.Username

	// Send registration success
	h.sendJSON(client, map[string]interface{}{
		"type": "registered",
		"data": user,
	})

	// Send current open orders
	if orders, err := h.orderService.GetAllOpenOrders(); err == nil {
		h.sendJSON(client, map[string]interface{}{
			"type": "orders_update",
			"data": orders,
		})
	}
}

func (h *Handler) handleChatMessage(client *ws.Client, msg *WSMessage) {
	user := h.clientUsers[client]
	if user == nil {
		h.sendError(client, "Please register first")
		return
	}

	content := h.extractMessageContent(msg.Data)
	if content == "" {
		h.sendError(client, "Invalid message format")
		return
	}

	// Save and broadcast message
	h.messageService.SaveMessage(content, user.Username)
	h.broadcastMessage(user.Username, content)
}

func (h *Handler) handleOrder(client *ws.Client, msg *WSMessage) {
	user := h.clientUsers[client]
	if user == nil {
		h.sendError(client, "Please register first")
		return
	}

	var data OrderData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		h.sendError(client, "Invalid order data")
		return
	}

	// Validate
	if err := h.validateOrderData(&data); err != nil {
		h.sendError(client, err.Error())
		return
	}

	orderType := model.OrderType(data.OrderType)
	order, err := h.orderService.CreateOrder(user.ID, user.Username, data.Symbol, orderType, data.Price, data.Quantity)
	if err != nil {
		h.sendError(client, "Failed to create order: "+err.Error())
		return
	}

	// Broadcast order to all clients
	h.hub.BroadcastJSON(map[string]interface{}{
		"type": "order_created",
		"data": order,
	})
}

func (h *Handler) handleCancelOrder(client *ws.Client, msg *WSMessage) {
	user := h.clientUsers[client]
	if user == nil {
		h.sendError(client, "Please register first")
		return
	}

	var data CancelOrderData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		h.sendError(client, "Invalid cancel order data")
		return
	}

	if err := h.orderService.CancelOrder(data.OrderID, user.ID); err != nil {
		h.sendError(client, "Failed to cancel order: "+err.Error())
		return
	}

	// Broadcast cancellation
	h.hub.BroadcastJSON(map[string]interface{}{
		"type": "order_cancelled",
		"data": map[string]interface{}{
			"order_id": data.OrderID,
		},
	})
}

// Helper methods

func (h *Handler) extractMessageContent(data json.RawMessage) string {
	var content string

	// Try to unmarshal as string first
	if err := json.Unmarshal(data, &content); err == nil {
		return content
	}

	// Try as object with content field
	var msgData map[string]string
	if err := json.Unmarshal(data, &msgData); err == nil {
		return msgData["content"]
	}

	// Try as plain string (backward compatibility)
	content = string(data)
	if len(content) > 2 && content[0] == '"' && content[len(content)-1] == '"' {
		return content[1 : len(content)-1]
	}

	return ""
}

func (h *Handler) validateOrderData(data *OrderData) error {
	orderType := model.OrderType(data.OrderType)
	if orderType != model.OrderTypeBid && orderType != model.OrderTypeAsk {
		return ErrInvalidOrderType
	}
	if data.Price <= 0 {
		return ErrInvalidPrice
	}
	if data.Quantity <= 0 {
		return ErrInvalidQuantity
	}
	if data.Symbol == "" {
		return ErrEmptySymbol
	}
	return nil
}

func (h *Handler) sendError(client *ws.Client, message string) {
	h.sendJSON(client, map[string]interface{}{
		"type": "error",
		"data": message,
	})
}

func (h *Handler) sendJSON(client *ws.Client, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}
	client.Send(jsonData)
}

func (h *Handler) broadcastMessage(username, content string) {
	h.hub.BroadcastJSON(map[string]interface{}{
		"type": "message",
		"data": map[string]string{
			"username": username,
			"content":  content,
		},
	})
}
