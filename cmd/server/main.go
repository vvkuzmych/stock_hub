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

var upgraderPG = gorillaWS.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

// WSMessagePG represents a WebSocket message
type WSMessagePG struct {
	Type     string          `json:"type"` // "register", "message", "order", "cancel_order"
	Data     json.RawMessage `json:"data"`
	Username string          `json:"username,omitempty"`
	UserID   int64           `json:"user_id,omitempty"`
}

// RegisterDataPG for user registration
type RegisterDataPG struct {
	Username string `json:"username"`
}

// OrderDataPG for stock orders
type OrderDataPG struct {
	Symbol    string  `json:"symbol"`
	OrderType string  `json:"order_type"` // "bid" or "ask"
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
}

// CancelOrderDataPG for cancelling orders
type CancelOrderDataPG struct {
	OrderID int64 `json:"order_id"`
}

type ServerContextPG struct {
	messageService *service.MessageServicePostgres
	userService    *service.UserService
	orderService   *service.StockOrderService
	hub            *ws.Hub
	clientUsers    map[*ws.Client]*model.User // Map clients to users
}

func main() {
	cfg := config.LoadPostgres()

	log.Printf("Starting Stock Hub with %s database", cfg.DatabaseType)
	log.Printf("DSN: %s", cfg.GetDSN())

	// Initialize message service with database type support
	messageService, err := service.NewMessageServicePostgres(
		cfg.GetDriver(),
		cfg.GetDSN(),
		cfg.GetMigrationsDir(),
	)
	if err != nil {
		log.Fatalf("Failed to initialize message service: %v", err)
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

	ctx := &ServerContextPG{
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
		handleWebSocketPG(ctx, w, r)
	})

	// Serve index.html for root and all other routes (SPA routing)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/ws" {
			return
		}
		// Serve index.html directly
		http.ServeFile(w, r, "static/index.html")
	})

	log.Printf("🚀 Server starting on :%s", cfg.ServerPort)
	log.Printf("🌐 Web interface: http://localhost:%s", cfg.ServerPort)
	log.Printf("📊 Database: %s", cfg.DatabaseType)
	if err := http.ListenAndServe(":"+cfg.ServerPort, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func handleWebSocketPG(ctx *ServerContextPG, w http.ResponseWriter, r *http.Request) {
	log.Printf("WebSocket connection attempt from %s", r.RemoteAddr)
	conn, err := upgraderPG.Upgrade(w, r, nil)
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

			var msg WSMessagePG
			if err := json.Unmarshal(messageBytes, &msg); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}

			handleMessagePG(ctx, client, &msg)
		}
	}()
}

func handleMessagePG(ctx *ServerContextPG, client *ws.Client, msg *WSMessagePG) {
	switch msg.Type {
	case "register":
		var data RegisterDataPG
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			sendErrorPG(client, "Invalid registration data")
			return
		}

		user, err := ctx.userService.RegisterUser(data.Username)
		if err != nil {
			sendErrorPG(client, "Failed to register user: "+err.Error())
			return
		}

		ctx.clientUsers[client] = user
		client.ID = data.Username

		// Send registration success
		response := map[string]interface{}{
			"type": "registered",
			"data": user,
		}
		sendJSONPG(client, response)

		// Send current open orders
		orders, err := ctx.orderService.GetAllOpenOrders()
		if err == nil {
			response = map[string]interface{}{
				"type": "orders_update",
				"data": orders,
			}
			sendJSONPG(client, response)
		}

	case "message":
		user := ctx.clientUsers[client]
		if user == nil {
			sendErrorPG(client, "Please register first")
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
			sendErrorPG(client, "Invalid message format")
			return
		}

		// Save and broadcast message
		ctx.messageService.SaveMessage(content, user.Username)
		broadcastMessagePG(ctx, user.Username, content)

	case "order":
		user := ctx.clientUsers[client]
		if user == nil {
			sendErrorPG(client, "Please register first")
			return
		}

		var data OrderDataPG
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			sendErrorPG(client, "Invalid order data")
			return
		}

		orderType := model.OrderType(data.OrderType)
		if orderType != model.OrderTypeBid && orderType != model.OrderTypeAsk {
			sendErrorPG(client, "Invalid order type")
			return
		}

		order, err := ctx.orderService.CreateOrder(user.ID, user.Username, data.Symbol, orderType, data.Price, data.Quantity)
		if err != nil {
			sendErrorPG(client, "Failed to create order: "+err.Error())
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
			sendErrorPG(client, "Please register first")
			return
		}

		var data CancelOrderDataPG
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			sendErrorPG(client, "Invalid cancel order data")
			return
		}

		if err := ctx.orderService.CancelOrder(data.OrderID, user.ID); err != nil {
			sendErrorPG(client, "Failed to cancel order: "+err.Error())
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

func sendErrorPG(client *ws.Client, message string) {
	response := map[string]interface{}{
		"type": "error",
		"data": message,
	}
	sendJSONPG(client, response)
}

func sendJSONPG(client *ws.Client, data interface{}) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}
	client.Send(jsonData)
}

func broadcastMessagePG(ctx *ServerContextPG, username, content string) {
	response := map[string]interface{}{
		"type": "message",
		"data": map[string]string{
			"username": username,
			"content":  content,
		},
	}
	ctx.hub.BroadcastJSON(response)
}
