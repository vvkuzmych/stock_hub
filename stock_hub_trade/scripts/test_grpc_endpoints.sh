#!/bin/bash

# Test all gRPC endpoints
# Usage: ./scripts/test_grpc_endpoints.sh

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}🧪 Testing gRPC Endpoints${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if server is running
echo -n "1️⃣  Checking server... "
if ! lsof -i :50051 > /dev/null 2>&1; then
    echo -e "${RED}❌ Server not running${NC}"
    echo ""
    echo "Start server first:"
    echo "  make run"
    exit 1
fi
echo -e "${GREEN}✅ Running on :50051${NC}"
echo ""

# List services
echo -n "2️⃣  Listing services... "
SERVICES=$(grpcurl -plaintext localhost:50051 list 2>&1)
if echo "$SERVICES" | grep -q "stock_hub.StockHubService"; then
    echo -e "${GREEN}✅${NC}"
    echo "   Available services:"
    echo "$SERVICES" | grep -v "reflection" | sed 's/^/   • /'
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$SERVICES"
    exit 1
fi
echo ""

# Test GetUser
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}📋 Test 1: GetUser${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Request: {\"user_id\": 1}"
echo ""
RESPONSE=$(grpcurl -plaintext \
  -d '{"user_id": 1}' \
  localhost:50051 \
  stock_hub.StockHubService/GetUser 2>&1)

if echo "$RESPONSE" | grep -q "username"; then
    echo -e "${GREEN}✅ Success!${NC}"
    echo "$RESPONSE" | jq '.'
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test GetOrder
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}📋 Test 2: GetOrder${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Request: {\"order_id\": 1}"
echo ""
RESPONSE=$(grpcurl -plaintext \
  -d '{"order_id": 1}' \
  localhost:50051 \
  stock_hub.StockHubService/GetOrder 2>&1)

if echo "$RESPONSE" | grep -q "order" || echo "$RESPONSE" | grep -q "NotFound"; then
    echo -e "${GREEN}✅ Success!${NC}"
    echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test GetOrdersByUser
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}📋 Test 3: GetOrdersByUser${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Request: {\"user_id\": 1, \"limit\": 10}"
echo ""
RESPONSE=$(grpcurl -plaintext \
  -d '{"user_id": 1, "limit": 10}' \
  localhost:50051 \
  stock_hub.StockHubService/GetOrdersByUser 2>&1)

if echo "$RESPONSE" | grep -q "orders" || echo "$RESPONSE" | grep -q "\[\]"; then
    echo -e "${GREEN}✅ Success!${NC}"
    echo "$RESPONSE" | jq '.'
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test GetOrders
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}📋 Test 4: GetOrders (all)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Request: {}"
echo ""
RESPONSE=$(grpcurl -plaintext \
  -d '{}' \
  localhost:50051 \
  stock_hub.StockHubService/GetOrders 2>&1)

if echo "$RESPONSE" | grep -q "orders" || echo "$RESPONSE" | grep -q "\[\]"; then
    echo -e "${GREEN}✅ Success!${NC}"
    echo "$RESPONSE" | jq '.'
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test GetMessages
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}📋 Test 5: GetMessages${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Request: {\"limit\": 10}"
echo ""
RESPONSE=$(grpcurl -plaintext \
  -d '{"limit": 10}' \
  localhost:50051 \
  stock_hub.StockHubService/GetMessages 2>&1)

if echo "$RESPONSE" | grep -q "messages" || echo "$RESPONSE" | grep -q "\[\]"; then
    echo -e "${GREEN}✅ Success!${NC}"
    echo "$RESPONSE" | jq '.'
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}🎉 All tests completed!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "💡 Try Web UI for interactive testing:"
echo "   ./scripts/start_grpcui.sh"
echo ""
