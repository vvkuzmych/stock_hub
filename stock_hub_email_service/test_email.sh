#!/bin/bash
# Quick Email Test - Перевірка чи працює відправка

set -e

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📧 Email Service Test"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. Check Mailhog
echo -n "1️⃣  Mailhog running... "
if curl -s http://localhost:8025 > /dev/null 2>&1; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${RED}❌${NC}"
    echo ""
    echo "   Start Mailhog: docker-compose up -d"
    exit 1
fi

# 2. Check Email Service
echo -n "2️⃣  Email service running... "
if curl -s http://localhost:8083/health | grep -q "ok"; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${RED}❌${NC}"
    echo ""
    echo "   Start service: make run"
    exit 1
fi

# 3. Check Database
echo -n "3️⃣  Database connected... "
if psql -U postgres -h localhost -d stock_hub -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${RED}❌${NC}"
    echo ""
    echo "   Start PostgreSQL: cd ../stock_hub_trade && docker-compose up -d"
    exit 1
fi

# 4. Get user from database and verify email
echo -n "4️⃣  Getting user from DB... "
USER_DATA=$(psql -U postgres -h localhost -d stock_hub -t -c "SELECT id, username, COALESCE(email, '') as email FROM users ORDER BY id LIMIT 1;" 2>/dev/null)

if [ -z "$USER_DATA" ]; then
    echo -e "${RED}❌ No users found${NC}"
    echo ""
    echo "   Creating test user with real email..."
    psql -U postgres -h localhost -d stock_hub -c "INSERT INTO users (username, email) VALUES ('test_user', 'test@stockhub.com');" > /dev/null 2>&1
    USER_DATA=$(psql -U postgres -h localhost -d stock_hub -t -c "SELECT id, username, email FROM users WHERE username = 'test_user';")
fi

USER_ID=$(echo "$USER_DATA" | awk '{print $1}' | tr -d ' ')
USER_NAME=$(echo "$USER_DATA" | awk '{print $2}' | tr -d ' ')
USER_EMAIL=$(echo "$USER_DATA" | awk '{print $3}' | tr -d ' ')

# Validate email format (must contain @)
if [ -z "$USER_EMAIL" ] || ! echo "$USER_EMAIL" | grep -q "@"; then
    echo -e "${YELLOW}⚠️  Invalid email '$USER_EMAIL', fixing...${NC}"
    USER_EMAIL="${USER_NAME}@stockhub.com"
    psql -U postgres -h localhost -d stock_hub -c "UPDATE users SET email = '$USER_EMAIL' WHERE id = $USER_ID;" > /dev/null 2>&1
fi

echo -e "${GREEN}✅ ID=$USER_ID, Email=$USER_EMAIL${NC}"

# 5. Clear Mailhog inbox
echo -n "5️⃣  Clearing inbox... "
curl -s -X DELETE http://localhost:8025/api/v1/messages > /dev/null 2>&1
echo -e "${GREEN}✅${NC}"

# 6. Send test email
echo -n "6️⃣  Sending test email to $USER_EMAIL... "
RESPONSE=$(curl -s -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d "{\"type\": \"welcome\", \"user_id\": $USER_ID}")

if echo "$RESPONSE" | grep -q "sent"; then
    echo -e "${GREEN}✅${NC}"
else
    echo -e "${RED}❌ $RESPONSE${NC}"
    exit 1
fi

# 6. Wait for email
echo -n "6️⃣  Checking Mailhog... "
sleep 2
TOTAL=$(curl -s http://localhost:8025/api/v2/messages | jq -r '.total // 0')

if [ "$TOTAL" -gt 0 ]; then
    echo -e "${GREEN}✅ Email received!${NC}"
    echo ""
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo -e "${GREEN}🎉 SUCCESS! Email system working!${NC}"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    echo ""
    echo "📧 Email Details:"
    curl -s http://localhost:8025/api/v2/messages | jq -r '.items[0] | "   From: \(.Content.Headers.From[0])\n   To: \(.Content.Headers.To[0])\n   Subject: \(.Content.Headers.Subject[0])"'
    echo ""
    echo "🌐 View in browser: http://localhost:8025"
    echo ""
    
    # Try to open browser
    if command -v open > /dev/null 2>&1; then
        open http://localhost:8025
    fi
    
    exit 0
else
    echo -e "${RED}❌ No email in Mailhog${NC}"
    echo ""
    echo "Troubleshooting:"
    echo "  - Check email service logs"
    echo "  - Verify SMTP_HOST=localhost, SMTP_PORT=1025 in .env"
    echo "  - Ensure MOCK_EMAIL=false"
    exit 1
fi
