package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"stock_hub_trade/api/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// APIGateway handles REST API requests and forwards them to gRPC
type APIGateway struct {
	grpcClient proto.StockHubServiceClient
}

// NewAPIGateway creates a new API gateway
func NewAPIGateway(grpcClient proto.StockHubServiceClient) *APIGateway {
	return &APIGateway{
		grpcClient: grpcClient,
	}
}

// GetUser handles GET /api/v1/users/:id
func (g *APIGateway) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	userIDStr := r.URL.Path[len("/api/v1/users/"):]
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		g.errorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Call gRPC service
	resp, err := g.grpcClient.GetUser(r.Context(), &proto.GetUserRequest{
		UserId: userID,
	})
	if err != nil {
		g.handleGRPCError(w, err)
		return
	}

	// Return JSON response
	g.jsonResponse(w, resp.User, http.StatusOK)
}

// GetOrder handles GET /api/v1/orders/:id
func (g *APIGateway) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Extract order ID from URL path
	orderIDStr := r.URL.Path[len("/api/v1/orders/"):]
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		g.errorResponse(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Call gRPC service
	resp, err := g.grpcClient.GetOrder(r.Context(), &proto.GetOrderRequest{
		OrderId: orderID,
	})
	if err != nil {
		g.handleGRPCError(w, err)
		return
	}

	// Return JSON response
	g.jsonResponse(w, resp.Order, http.StatusOK)
}

// GetOrdersByUser handles GET /api/v1/users/:id/orders
func (g *APIGateway) GetOrdersByUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		g.errorResponse(w, "user_id parameter is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		g.errorResponse(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get limit parameter (default: 50)
	limit := int32(50)
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	// Call gRPC service
	resp, err := g.grpcClient.GetOrdersByUser(r.Context(), &proto.GetOrdersByUserRequest{
		UserId: userID,
		Limit:  limit,
	})
	if err != nil {
		g.handleGRPCError(w, err)
		return
	}

	// Return JSON response
	g.jsonResponse(w, resp.Orders, http.StatusOK)
}

// GetOrders handles GET /api/v1/orders
func (g *APIGateway) GetOrders(w http.ResponseWriter, r *http.Request) {
	// Get optional symbol filter
	symbol := r.URL.Query().Get("symbol")

	// Call gRPC service
	resp, err := g.grpcClient.GetOrders(r.Context(), &proto.GetOrdersRequest{
		Symbol: symbol,
	})
	if err != nil {
		g.handleGRPCError(w, err)
		return
	}

	// Return JSON response
	g.jsonResponse(w, resp.Orders, http.StatusOK)
}

// GetMessages handles GET /api/v1/messages
func (g *APIGateway) GetMessages(w http.ResponseWriter, r *http.Request) {
	// Get limit parameter (default: 100)
	limit := int32(100)
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.ParseInt(limitStr, 10, 32); err == nil {
			limit = int32(l)
		}
	}

	// Call gRPC service
	resp, err := g.grpcClient.GetMessages(r.Context(), &proto.GetMessagesRequest{
		Limit: limit,
	})
	if err != nil {
		g.handleGRPCError(w, err)
		return
	}

	// Return JSON response
	g.jsonResponse(w, resp.Messages, http.StatusOK)
}

// Helper methods

// jsonResponse sends JSON response
func (g *APIGateway) jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error but don't send another response
		return
	}
}

// errorResponse sends error JSON response
func (g *APIGateway) errorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"error":  message,
		"status": statusCode,
	}

	json.NewEncoder(w).Encode(response)
}

// handleGRPCError converts gRPC error to HTTP error
func (g *APIGateway) handleGRPCError(w http.ResponseWriter, err error) {
	// Get gRPC status
	st, ok := status.FromError(err)
	if !ok {
		g.errorResponse(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Map gRPC codes to HTTP status codes
	var httpStatus int
	switch st.Code() {
	case codes.OK:
		httpStatus = http.StatusOK
	case codes.InvalidArgument:
		httpStatus = http.StatusBadRequest
	case codes.NotFound:
		httpStatus = http.StatusNotFound
	case codes.AlreadyExists:
		httpStatus = http.StatusConflict
	case codes.PermissionDenied:
		httpStatus = http.StatusForbidden
	case codes.Unauthenticated:
		httpStatus = http.StatusUnauthorized
	case codes.ResourceExhausted:
		httpStatus = http.StatusTooManyRequests
	case codes.Unimplemented:
		httpStatus = http.StatusNotImplemented
	case codes.Unavailable:
		httpStatus = http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		httpStatus = http.StatusGatewayTimeout
	default:
		httpStatus = http.StatusInternalServerError
	}

	g.errorResponse(w, st.Message(), httpStatus)
}
