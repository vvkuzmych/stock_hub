package grpc

import (
	"context"
	"fmt"
	"net"

	"stock_hub_trade/api/proto"
	"stock_hub_trade/internal/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

// Server implements the StockHubService gRPC service
type Server struct {
	proto.UnimplementedStockHubServiceServer
	userService    *service.UserService
	orderService   *service.StockOrderService
	messageService *service.MessageService
}

// NewServer creates a new gRPC server
func NewServer(
	userService *service.UserService,
	orderService *service.StockOrderService,
	messageService *service.MessageService,
) *Server {
	return &Server{
		userService:    userService,
		orderService:   orderService,
		messageService: messageService,
	}
}

// GetUser retrieves a user by ID
func (s *Server) GetUser(ctx context.Context, req *proto.GetUserRequest) (*proto.GetUserResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

	user, err := s.userService.GetUserByID(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &proto.GetUserResponse{
		User: ToProtoUser(user),
	}, nil
}

// GetUsers retrieves all users (with pagination)
func (s *Server) GetUsers(ctx context.Context, req *proto.GetUsersRequest) (*proto.GetUsersResponse, error) {
	// For now, return empty list as we don't have GetUsers in service
	// TODO: Implement GetUsers in UserService
	return &proto.GetUsersResponse{
		Users: []*proto.User{},
	}, nil
}

// GetOrder retrieves an order by ID
func (s *Server) GetOrder(ctx context.Context, req *proto.GetOrderRequest) (*proto.GetOrderResponse, error) {
	if req.OrderId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "order_id must be positive")
	}

	order, err := s.orderService.GetOrderByID(req.OrderId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "order not found: %v", err)
	}

	return &proto.GetOrderResponse{
		Order: ToProtoStockOrder(order),
	}, nil
}

// GetOrders retrieves orders with filters
func (s *Server) GetOrders(ctx context.Context, req *proto.GetOrdersRequest) (*proto.GetOrdersResponse, error) {
	var orders []*proto.StockOrder

	if req.Symbol != "" {
		// Get orders for specific symbol
		domainOrders, err := s.orderService.GetOpenOrders(req.Symbol)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get orders: %v", err)
		}
		orders = ToProtoStockOrders(domainOrders)
	} else {
		// Get all open orders
		domainOrders, err := s.orderService.GetAllOpenOrders()
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to get orders: %v", err)
		}
		orders = ToProtoStockOrders(domainOrders)
	}

	return &proto.GetOrdersResponse{
		Orders: orders,
	}, nil
}

// GetOrdersByUser retrieves orders for a specific user
func (s *Server) GetOrdersByUser(ctx context.Context, req *proto.GetOrdersByUserRequest) (*proto.GetOrdersByUserResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id must be positive")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 100
	}

	orders, err := s.orderService.GetOrdersByUserID(req.UserId, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get orders: %v", err)
	}

	return &proto.GetOrdersByUserResponse{
		Orders: ToProtoStockOrders(orders),
	}, nil
}

// GetMessages retrieves messages
func (s *Server) GetMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 100 // Default limit
	}

	messages, err := s.messageService.GetMessages(limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get messages: %v", err)
	}

	return &proto.GetMessagesResponse{
		Messages: ToProtoMessages(messages),
	}, nil
}

// GetMessagesByClient retrieves messages for a specific client
func (s *Server) GetMessagesByClient(ctx context.Context, req *proto.GetMessagesByClientRequest) (*proto.GetMessagesByClientResponse, error) {
	if req.ClientId == "" {
		return nil, status.Error(codes.InvalidArgument, "client_id cannot be empty")
	}

	limit := int(req.Limit)
	if limit <= 0 {
		limit = 100 // Default limit
	}

	messages, err := s.messageService.GetMessagesByClient(req.ClientId, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get messages: %v", err)
	}

	return &proto.GetMessagesByClientResponse{
		Messages: ToProtoMessages(messages),
	}, nil
}

// RegisterServer registers the gRPC server with a gRPC server instance
func RegisterServer(
	grpcServer *grpc.Server,
	userService *service.UserService,
	orderService *service.StockOrderService,
	messageService *service.MessageService,
) {
	server := NewServer(userService, orderService, messageService)
	proto.RegisterStockHubServiceServer(grpcServer, server)
}

// NewGRPCServer creates a new gRPC server instance with the StockHub service
func NewGRPCServer(
	userService *service.UserService,
	orderService *service.StockOrderService,
	messageService *service.MessageService,
) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	RegisterServer(grpcServer, userService, orderService, messageService)

	// Register reflection service for development/testing
	// Allows tools like Kreya and grpcurl to discover methods automatically
	reflection.Register(grpcServer)

	return grpcServer, nil
}

// ServeGRPC starts the gRPC server on the specified port
func ServeGRPC(
	port string,
	userService *service.UserService,
	orderService *service.StockOrderService,
	messageService *service.MessageService,
) error {
	grpcServer, err := NewGRPCServer(userService, orderService, messageService)
	if err != nil {
		return fmt.Errorf("failed to create gRPC server: %w", err)
	}

	// Listen on TCP port
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", port, err)
	}

	fmt.Printf("🚀 gRPC server listening on :%s\n", port)
	return grpcServer.Serve(lis)
}
