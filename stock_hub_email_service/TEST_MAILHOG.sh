#!/bin/bash
# Test Mailhog Integration

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🧪 Testing Mailhog Integration"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 1. Check Mailhog
echo "1️⃣ Checking Mailhog..."
if docker ps | grep -q mailhog; then
    echo "   ✅ Mailhog is running"
else
    echo "   ❌ Mailhog is NOT running"
    echo "   Run: docker-compose up -d"
    exit 1
fi
echo ""

# 2. Check Mailhog Web UI
echo "2️⃣ Checking Mailhog Web UI..."
if curl -s http://localhost:8025 > /dev/null 2>&1; then
    echo "   ✅ Web UI accessible: http://localhost:8025"
else
    echo "   ❌ Web UI not accessible"
    exit 1
fi
echo ""

# 3. Check Email Service
echo "3️⃣ Checking Email Service..."
if lsof -i :8083 > /dev/null 2>&1; then
    echo "   ✅ Email service is running on port 8083"
else
    echo "   ❌ Email service is NOT running"
    echo "   Run: make run"
    exit 1
fi
echo ""

# 4. Check Database
echo "4️⃣ Checking Database..."
if psql -U postgres -h localhost -d stock_hub -c "SELECT COUNT(*) FROM users;" > /dev/null 2>&1; then
    USER_COUNT=$(psql -U postgres -h localhost -d stock_hub -t -c "SELECT COUNT(*) FROM users;")
    echo "   ✅ Database connected: $USER_COUNT users"
else
    echo "   ❌ Database connection failed"
    exit 1
fi
echo ""

# 5. Clear old emails
echo "5️⃣ Clearing old emails..."
curl -s -X DELETE http://localhost:8025/api/v1/messages > /dev/null 2>&1
echo "   ✅ Mailhog inbox cleared"
echo ""

# 6. Send test email
echo "6️⃣ Sending test email..."
RESPONSE=$(curl -s -X POST http://localhost:8083/send \
  -H "Content-Type: application/json" \
  -d '{"type": "welcome", "user_id": 1}')

if echo "$RESPONSE" | grep -q "sent"; then
    echo "   ✅ Email sent: $RESPONSE"
else
    echo "   ❌ Failed to send email: $RESPONSE"
    exit 1
fi
echo ""

# 7. Wait and check Mailhog
echo "7️⃣ Checking Mailhog inbox..."
sleep 2

TOTAL=$(curl -s http://localhost:8025/api/v2/messages | jq -r '.total // 0')
echo "   📬 Emails in inbox: $TOTAL"

if [ "$TOTAL" -gt 0 ]; then
    echo "   ✅ Email received in Mailhog!"
    echo ""
    echo "   📧 Email Details:"
    curl -s http://localhost:8025/api/v2/messages | jq -r '.items[0] | "      From: \(.Content.Headers.From[0] // "N/A")\n      To: \(.Content.Headers.To[0] // "N/A")\n      Subject: \(.Content.Headers.Subject[0] // "N/A")"'
else
    echo "   ⚠️  No emails in Mailhog inbox"
    echo "   Check email service logs for errors"
fi
echo ""

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🌐 Open in browser:"
echo "   http://localhost:8025"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# 8. Try to open browser
if command -v open > /dev/null 2>&1; then
    open http://localhost:8025
elif command -v xdg-open > /dev/null 2>&1; then
    xdg-open http://localhost:8025
else
    echo "   Open manually: http://localhost:8025"
fi

echo "✅ Test complete!"
