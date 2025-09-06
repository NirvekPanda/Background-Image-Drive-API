#!/bin/bash

# Simple Health Check Test
# Tests if the API is currently live and responding

echo "🏥 API Health Check Test"
echo "========================"

# Load environment variables from test.env if it exists
if [ -f "testing/test.env" ]; then
    echo "📋 Loading test configuration from testing/test.env"
    export $(grep -v '^#' testing/test.env | xargs)
fi

# Configuration
API_BASE_URL="${API_BASE_URL:-https://background-image-drive-api-189526192204.us-west1.run.app}"

echo "Testing API at: $API_BASE_URL"
echo ""

# Test 1: Basic health endpoint
echo "1️⃣ Testing /health endpoint..."
HEALTH_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$API_BASE_URL/health")
HTTP_STATUS=$(echo $HEALTH_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
HEALTH_BODY=$(echo $HEALTH_RESPONSE | sed -e 's/HTTPSTATUS:.*//g')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ Health check passed (Status: $HTTP_STATUS)"
    echo "   Response: $HEALTH_BODY"
else
    echo "❌ Health check failed (Status: $HTTP_STATUS)"
    echo "   Response: $HEALTH_BODY"
    exit 1
fi

echo ""

# Test 2: Check if API is responsive
echo "2️⃣ Testing API responsiveness..."
RESPONSE_TIME=$(curl -s -w "%{time_total}" -o /dev/null "$API_BASE_URL/health")
echo "✅ API response time: ${RESPONSE_TIME}s"

echo ""

# Test 3: Check CORS headers
echo "3️⃣ Testing CORS headers..."
CORS_RESPONSE=$(curl -s -I -H "Origin: http://localhost:3000" "$API_BASE_URL/health")
if echo "$CORS_RESPONSE" | grep -q "Access-Control-Allow-Origin"; then
    echo "✅ CORS headers present"
else
    echo "⚠️  CORS headers not found"
fi

echo ""
echo "🎉 Health check completed successfully!"
echo "API is live and ready for testing."
