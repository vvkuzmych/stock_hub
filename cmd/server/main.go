package main

import (
	"encoding/json"
	"log"
	"net/http"

	"stock_hub/internal/config"
	"stock_hub/internal/model"
	"stock_hub/internal/service"
	ws "stock_hub/pkg/websocket"

	gorillaWS "github.com/gorilla/websocket"
)

var upgrader = gorillaWS.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// WSMessage represents a WebSocket message
type WSMessage struct {
	Type     string          `json:"type"` // "register", "message", "order", "cancel_order"
	Data     json.RawMessage `json:"data"`
	Username string          `json:"username,omitempty"`
	UserID   int64           `json:"user_id,omitempty"`
}

// RegisterData for user registration
type RegisterData struct {
	Username string `json:"username"`
}

// OrderData for stock orders
type OrderData struct {
	Symbol    string  `json:"symbol"`
	OrderType string  `json:"order_type"` // "bid" or "ask"
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

// CancelOrderData for cancelling orders
type CancelOrderData struct {
	OrderID int64 `json:"order_id"`
}

type ServerContext struct {
	messageService *service.MessageService
	userService    *service.UserService
	orderService   *service.StockOrderService
	hub            *ws.Hub
	clientUsers    map[*ws.Client]*model.User // Map clients to users
}

func main() {
	cfg := config.Load()

	log.Printf("🚀 Starting Stock Hub with %s database", cfg.DatabaseType)
	log.Printf("📊 DSN: %s", cfg.GetSafeDSN())

	// Initialize message service with database type support
	messageService, err := service.NewMessageService(
		cfg.GetDriver(),
		cfg.GetDSN(),
		cfg.GetMigrationsDir(),
	)
	if err != nil {
		log.Fatalf("❌ Failed to initialize message service: %v", err)
	}
	defer messageService.Close()

	log.Println("✅ Database connected and migrations applied")

	// Get database connection from message service
	db := messageService.GetDB()
	userService := service.NewUserServiceWithDriver(db, cfg.GetDriver())
	orderService := service.NewStockOrderServiceWithDriver(db, cfg.GetDriver())

	// Create WebSocket hub with message logging
	hub := ws.NewHub(messageService)
	go hub.Run()

	ctx := &ServerContext{
		messageService: messageService,
		userService:    userService,
		orderService:   orderService,
		hub:            hub,
		clientUsers:    make(map[*ws.Client]*model.User),
	}

	mux := http.NewServeMux()

	// Serve static assets
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("static/assets"))))

	// Register WebSocket handler
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(ctx, w, r)
	})

	// Serve index.html for root and all other routes (SPA routing)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws" {
			return
		}
		http.ServeFile(w, r, "static/index.html")
	})

	log.Printf("🌐 Web interface: http://localhost:%s", cfg.ServerPort)
	log.Printf("🔌 WebSocket endpoint: ws://localhost:%s/ws", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func handleWebSocket(ctx *ServerContext, w http.ResponseWriter, r *http.Request) {
	log.Printf("WebSocket connection attempt from %s", r.RemoteAddr)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}

	// Create client and register with hub
	client := ws.NewClient(ctx.hub, conn, "")
	ctx.hub.Register(client)

	// Start write pump
	go client.WritePump()

	// Handle messages in a separate goroutine
	// Note: We implement a custom read loop here instead of using Client.ReadPump()
	// because we need to parse structured WSMessage types (register, message, order, etc.)
	// and route them to different handlers, rather than simple broadcast
	go func() {
		defer func() {
			delete(ctx.clientUsers, client)
			ctx.hub.Unregister(client)
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

			handleMessage(ctx, client, &msg)
		}
	}()
}

func handleMessage(ctx *ServerContext, client *ws.Client, msg *WSMessage) {
	switch msg.Type {
	case "register":
		var data RegisterData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			sendError(client, "Invalid registration data")
			return
		}

		user, err := ctx.userService.RegisterUser(data.Username)
		if err != nil {
			sendError(client, "Failed to register user: "+err.Error())
			return
		}

		ctx.clientUsers[client] = user
		client.ID = data.Username

		// Send registration success
		response := map[string]interface{}{
			"type": "registered",
			"data": user,
		}
		sendJSON(client, response)

		// Send current open orders
		orders, err := ctx.orderService.GetAllOpenOrders()
		if err == nil {
			response = map[string]interface{}{
				"type": "orders_update",
				"data": orders,
			}
			sendJSON(client, response)
		}

	case "message":
		user := ctx.clientUsers[client]
		if user == nil {
			sendError(client, "Please register first")
			return
		}

		var content string
		// Try to unmarshal as string first
		if err := json.Unmarshal(msg.Data, &content); err != nil {
			// Try as object with content field
			var msgData map[string]string
			if err := json.Unmarshal(msg.Data, &msgData); err == nil {
				content = msgData["content"]
			} else {
				// Try as plain string (backward compatibility)
				content = string(msg.Data)
				if len(content) > 2 && content[0] == '"' && content[len(content)-1] == '"' {
					content = content[1 : len(content)-1]
				}
			}
		}

		if content == "" {
			sendError(client, "Invalid message format")
			return
		}

		// Save and broadcast message
		ctx.messageService.SaveMessage(content, user.Username)
		broadcastMessage(ctx, user.Username, content)

	case "order":
		user := ctx.clientUsers[client]
		if user == nil {
			sendError(client, "Please register first")
			return
		}

		var data OrderData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			sendError(client, "Invalid order data")
			return
		}

		orderType := model.OrderType(data.OrderType)
		if orderType != model.OrderTypeBid && orderType != model.OrderTypeAsk {
			sendError(client, "Invalid order type")
			return
		}

		order, err := ctx.orderService.CreateOrder(user.ID, user.Username, data.Symbol, orderType, data.Price, data.Quantity)
		if err != nil {
			sendError(client, "Failed to create order: "+err.Error())
			return
		}

		// Broadcast order to all clients
		response := map[string]interface{}{
			"type": "order_created",
			"data": order,
		}
		ctx.hub.BroadcastJSON(response)

	case "cancel_order":
		user := ctx.clientUsers[client]
		if user == nil {
			sendError(client, "Please register first")
			return
		}

		var data CancelOrderData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			sendError(client, "Invalid cancel order data")
			return
		}

		if err := ctx.orderService.CancelOrder(data.OrderID, user.ID); err != nil {
			sendError(client, "Failed to cancel order: "+err.Error())
			return
		}

		// Broadcast cancellation
		response := map[string]interface{}{
			"type": "order_cancelled",
			"data": map[string]interface{}{
				"order_id": data.OrderID,
			},
		}
		ctx.hub.BroadcastJSON(response)
	}
}

func sendError(client *ws.Client, message string) {
	response := map[string]interface{}{
		"type": "error",
		"data": message,
	}
	sendJSON(client, response)
}

func sendJSON(client *ws.Client, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}
	client.Send(jsonData)
}

func broadcastMessage(ctx *ServerContext, username, content string) {
	response := map[string]interface{}{
		"type": "message",
		"data": map[string]string{
			"username": username,
			"content":  content,
		},
	}
	ctx.hub.BroadcastJSON(response)
}
