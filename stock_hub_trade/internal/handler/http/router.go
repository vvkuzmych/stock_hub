package http

import (
	"net/http"

	wsHandler "stock_hub_trade/internal/handler/websocket"
	"stock_hub_trade/internal/middleware"
)

// Router holds routing configuration
type Router struct {
	mux        *http.ServeMux
	httpH      *Handler
	wsH        *wsHandler.Handler
	apiGateway *APIGateway
	middleware *middleware.Manager
}

// NewRouter creates a new router
func NewRouter(httpH *Handler, wsH *wsHandler.Handler, apiGateway *APIGateway, mw *middleware.Manager) *Router {
	return &Router{
		mux:        http.NewServeMux(),
		httpH:      httpH,
		wsH:        wsH,
		apiGateway: apiGateway,
		middleware: mw,
	}
}

// Setup configures all routes
func (r *Router) Setup() http.Handler {
	// Health check (no middleware)
	r.mux.HandleFunc("/health", r.httpH.HealthCheck)

	// Static assets
	r.mux.Handle("/assets/", r.httpH.ServeAssets())

	// WebSocket endpoint (with logging)
	r.mux.HandleFunc("/ws", r.middleware.Chain(
		r.wsH.HandleConnection,
		r.middleware.Logging,
		r.middleware.Recovery,
	))

	// REST API Gateway (REST → gRPC)
	r.setupAPIRoutes()

	// SPA routing - serve index.html for all other routes
	r.mux.HandleFunc("/", r.middleware.Chain(
		r.spaHandler,
		r.middleware.Logging,
		r.middleware.Recovery,
	))

	// Apply global middleware
	handler := r.middleware.CORS(r.mux)

	return handler
}

// setupAPIRoutes configures REST API routes (Gateway to gRPC)
func (r *Router) setupAPIRoutes() {
	// Users
	r.mux.HandleFunc("/api/v1/users/", r.middleware.Chain(
		r.apiGateway.GetUser,
		r.middleware.Logging,
		r.middleware.Recovery,
	))

	// Orders
	r.mux.HandleFunc("/api/v1/orders/", r.middleware.Chain(
		r.handleOrderRoutes,
		r.middleware.Logging,
		r.middleware.Recovery,
	))

	// Messages
	r.mux.HandleFunc("/api/v1/messages", r.middleware.Chain(
		r.apiGateway.GetMessages,
		r.middleware.Logging,
		r.middleware.Recovery,
	))
}

// handleOrderRoutes routes order-related requests
func (r *Router) handleOrderRoutes(w http.ResponseWriter, req *http.Request) {
	// Check if it's a specific order by ID
	path := req.URL.Path
	if len(path) > len("/api/v1/orders/") {
		// GET /api/v1/orders/:id
		r.apiGateway.GetOrder(w, req)
		return
	}

	// Check for user_id query parameter
	if req.URL.Query().Get("user_id") != "" {
		// GET /api/v1/orders?user_id=123
		r.apiGateway.GetOrdersByUser(w, req)
		return
	}

	// GET /api/v1/orders (all orders)
	r.apiGateway.GetOrders(w, req)
}

// spaHandler handles SPA routing
func (r *Router) spaHandler(w http.ResponseWriter, req *http.Request) {
	// Don't serve index.html for /ws endpoint
	if req.URL.Path == "/ws" {
		http.NotFound(w, req)
		return
	}
	r.httpH.ServeIndex(w, req)
}

// Handler returns the configured HTTP handler
func (r *Router) Handler() http.Handler {
	return r.Setup()
}
