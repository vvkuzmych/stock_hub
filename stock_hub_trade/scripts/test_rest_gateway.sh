#!/bin/bash

# Test REST API Gateway
# Usage: ./scripts/test_rest_gateway.sh

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

BASE_URL="http://localhost:8082/api/v1"

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}🧪 Testing REST API Gateway${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if server is running
echo -n "1️⃣  Checking server... "
if ! curl -s http://localhost:8082/health > /dev/null 2>&1; then
    echo -e "${RED}❌ Server not running${NC}"
    echo ""
    echo "Start server first:"
    echo "  cd /Users/vkuzm/GolandProjects/stock_hub/stock_hub_trade"
    echo "  make run"
    exit 1
fi
echo -e "${GREEN}✅ Running${NC}"
echo ""

# Test Health Check
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}Test 1: Health Check${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "GET http://localhost:8082/health"
echo ""
RESPONSE=$(curl -s http://localhost:8082/health)
echo "$RESPONSE" | jq .
echo ""

# Test Get User
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}Test 2: Get User${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "GET ${BASE_URL}/users/1"
echo ""
RESPONSE=$(curl -s ${BASE_URL}/users/1)
if echo "$RESPONSE" | jq -e '.id' > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Success!${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test Get Order
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}Test 3: Get Order${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "GET ${BASE_URL}/orders/1"
echo ""
RESPONSE=$(curl -s ${BASE_URL}/orders/1)
if echo "$RESPONSE" | jq -e '.id' > /dev/null 2>&1; then
    echo -e "${GREEN}✅ Success!${NC}"
    echo "$RESPONSE" | jq .
elif echo "$RESPONSE" | jq -e '.error' > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  No order found (expected if DB is empty)${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test Get Orders by User
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}Test 4: Get Orders by User${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "GET ${BASE_URL}/orders?user_id=1&limit=5"
echo ""
RESPONSE=$(curl -s "${BASE_URL}/orders?user_id=1&limit=5")
if echo "$RESPONSE" | jq -e 'type == "array"' > /dev/null 2>&1; then
    COUNT=$(echo "$RESPONSE" | jq 'length')
    echo -e "${GREEN}✅ Success! Found ${COUNT} orders${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test Get All Orders
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}Test 5: Get All Orders${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "GET ${BASE_URL}/orders"
echo ""
RESPONSE=$(curl -s ${BASE_URL}/orders)
if echo "$RESPONSE" | jq -e 'type == "array"' > /dev/null 2>&1; then
    COUNT=$(echo "$RESPONSE" | jq 'length')
    echo -e "${GREEN}✅ Success! Found ${COUNT} orders${NC}"
    echo "$RESPONSE" | jq 'if length > 3 then .[0:3] + ["..."] else . end'
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test Get Messages
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}Test 6: Get Messages${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "GET ${BASE_URL}/messages?limit=10"
echo ""
RESPONSE=$(curl -s "${BASE_URL}/messages?limit=10")
if echo "$RESPONSE" | jq -e 'type == "array"' > /dev/null 2>&1; then
    COUNT=$(echo "$RESPONSE" | jq 'length')
    echo -e "${GREEN}✅ Success! Found ${COUNT} messages${NC}"
    echo "$RESPONSE" | jq 'if length > 3 then .[0:3] + ["..."] else . end'
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$RESPONSE"
fi
echo ""

# Test 404 Error
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${BLUE}Test 7: Error Handling (404)${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "GET ${BASE_URL}/users/99999"
echo ""
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" ${BASE_URL}/users/99999)
RESPONSE=$(curl -s ${BASE_URL}/users/99999)
if [ "$HTTP_CODE" == "404" ]; then
    echo -e "${GREEN}✅ Error handling works! (404)${NC}"
    echo "$RESPONSE" | jq .
else
    echo -e "${YELLOW}⚠️  Expected 404, got ${HTTP_CODE}${NC}"
    echo "$RESPONSE"
fi
echo ""

# Summary
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}🎉 All tests completed!${NC}"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""
echo "📋 Available endpoints:"
echo "  GET  /api/v1/users/:id"
echo "  GET  /api/v1/orders/:id"
echo "  GET  /api/v1/orders?user_id=:id"
echo "  GET  /api/v1/orders"
echo "  GET  /api/v1/messages?limit=:n"
echo ""
echo "📚 Documentation:"
echo "  REST_GATEWAY_QUICKSTART.md"
echo "  docs/REST_API_GATEWAY.md"
echo ""
