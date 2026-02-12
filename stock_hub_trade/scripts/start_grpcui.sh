#!/bin/bash

# Start gRPC Web UI for Stock Hub
# Usage: ./scripts/start_grpcui.sh

set -e

echo "🚀 Starting gRPC Web UI for Stock Hub..."
echo ""

# Check if grpcui is installed
if ! command -v grpcui &> /dev/null; then
    echo "❌ grpcui not found!"
    echo ""
    echo "📦 Installing grpcui..."
    go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
    echo "✅ grpcui installed!"
    echo ""
fi

# Check if server is running
if ! lsof -i :50051 > /dev/null 2>&1; then
    echo "❌ gRPC server not running on :50051"
    echo ""
    echo "Start server first:"
    echo "  cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade"
    echo "  make run"
    exit 1
fi

echo "✅ gRPC server is running on :50051"
echo ""
echo "🌐 Opening Web UI..."
echo "   URL: http://localhost:8080 (will auto-open)"
echo ""
echo "📋 Available methods:"
echo "   • GetUser"
echo "   • GetUsers"
echo "   • GetOrder"
echo "   • GetOrders"
echo "   • GetOrdersByUser"
echo "   • GetMessages"
echo "   • GetMessagesByClient"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Start grpcui (will open browser automatically)
grpcui -plaintext localhost:50051
