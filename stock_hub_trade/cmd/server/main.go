package main

import (
	"log"
	"net/http"

	"stock_hub_trade/internal/config"
	grpcServer "stock_hub_trade/internal/grpc"
	httpHandler "stock_hub_trade/internal/handler/http"
	wsHandler "stock_hub_trade/internal/handler/websocket"
	"stock_hub_trade/internal/middleware"
	"stock_hub_trade/internal/service"
	ws "stock_hub_trade/pkg/websocket"
)

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

	// Start gRPC server in separate goroutine
	go func() {
		log.Printf("🚀 Starting gRPC server on port %s", cfg.GRPCPort)
		if err := grpcServer.ServeGRPC(cfg.GRPCPort, userService, orderService, messageService); err != nil {
			log.Fatalf("❌ gRPC server failed: %v", err)
		}
	}()

	// Create WebSocket hub with message logging
	hub := ws.NewHub(messageService)
	go hub.Run()

	// Create gRPC client for REST gateway
	grpcClient, err := httpHandler.NewGRPCClient(cfg.GRPCPort)
	if err != nil {
		log.Fatalf("❌ Failed to create gRPC client: %v", err)
	}
	defer grpcClient.Close()

	// Create handlers
	httpH := httpHandler.NewHandler()
	wsH := wsHandler.NewHandler(messageService, userService, orderService, hub)
	apiGateway := httpHandler.NewAPIGateway(grpcClient.Client())

	// Create middleware manager
	mw := middleware.NewManager()

	// Setup router
	router := httpHandler.NewRouter(httpH, wsH, apiGateway, mw)

	// Start HTTP server
	log.Printf("🌐 Web interface: http://localhost:%s", cfg.ServerPort)
	log.Printf("🔌 WebSocket endpoint: ws://localhost:%s/ws", cfg.ServerPort)
	log.Printf("🏥 Health check: http://localhost:%s/health", cfg.ServerPort)
	log.Printf("🚪 REST API Gateway: http://localhost:%s/api/v1/", cfg.ServerPort)

	if err := http.ListenAndServe(":"+cfg.ServerPort, router.Handler()); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
