package http

import (
	"fmt"

	"stock_hub_trade/api/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClient wraps gRPC client connection
type GRPCClient struct {
	conn   *grpc.ClientConn
	client proto.StockHubServiceClient
}

// NewGRPCClient creates internal gRPC client for REST gateway
func NewGRPCClient(port string) (*GRPCClient, error) {
	// Connect to internal gRPC server
	conn, err := grpc.Dial(
		"localhost:"+port,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := proto.NewStockHubServiceClient(conn)

	return &GRPCClient{
		conn:   conn,
		client: client,
	}, nil
}

// Client returns the gRPC service client
func (c *GRPCClient) Client() proto.StockHubServiceClient {
	return c.client
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
