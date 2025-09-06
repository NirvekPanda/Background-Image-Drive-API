#!/bin/bash

# Simple Location Test
# Tests location endpoints and extracts location metadata

echo "📍 Location Metadata Test"
echo "=========================="

# Load environment variables from test.env if it exists
if [ -f "testing/test.env" ]; then
    echo "📋 Loading test configuration from testing/test.env"
    export $(grep -v '^#' testing/test.env | xargs)
fi

# Configuration
API_BASE_URL="${API_BASE_URL:-https://background-image-drive-api-189526192204.us-west1.run.app}"

echo "API Base URL: $API_BASE_URL"
echo ""

# Test 1: Get location from coordinates
echo "1️⃣ Testing location from coordinates..."
LAT="${TEST_LATITUDE:-52.3676}"
LNG="${TEST_LONGITUDE:-4.9041}"

COORDS_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$API_BASE_URL/api/v1/location/coords?lat=$LAT&lng=$LNG")
HTTP_STATUS=$(echo $COORDS_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
COORDS_BODY=$(echo $COORDS_RESPONSE | sed -e 's/HTTPSTATUS:.*//g')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ Location from coordinates successful (Status: $HTTP_STATUS)"
    echo "   Coordinates: $LAT, $LNG"
    
    # Extract location metadata
    LOCATION_NAME=$(echo "$COORDS_BODY" | grep -o '"name":"[^"]*"' | cut -d'"' -f4)
    COUNTRY=$(echo "$COORDS_BODY" | grep -o '"country":"[^"]*"' | cut -d'"' -f4)
    CITY=$(echo "$COORDS_BODY" | grep -o '"city":"[^"]*"' | cut -d'"' -f4)
    
    echo "   Location: $LOCATION_NAME"
    echo "   City: $CITY"
    echo "   Country: $COUNTRY"
else
    echo "❌ Location from coordinates failed (Status: $HTTP_STATUS)"
    echo "   Response: $COORDS_BODY"
    exit 1
fi

echo ""

# Test 2: Get location from name
echo "2️⃣ Testing location from name..."
LOCATION_NAME="${TEST_LOCATION_NAME:-Amsterdam}"

NAME_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$API_BASE_URL/api/v1/location/name?name=$LOCATION_NAME")
HTTP_STATUS=$(echo $NAME_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
NAME_BODY=$(echo $NAME_RESPONSE | sed -e 's/HTTPSTATUS:.*//g')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ Location from name successful (Status: $HTTP_STATUS)"
    echo "   Search term: $LOCATION_NAME"
    
    # Extract location metadata
    FOUND_NAME=$(echo "$NAME_BODY" | grep -o '"name":"[^"]*"' | cut -d'"' -f4)
    FOUND_COUNTRY=$(echo "$NAME_BODY" | grep -o '"country":"[^"]*"' | cut -d'"' -f4)
    FOUND_CITY=$(echo "$NAME_BODY" | grep -o '"city":"[^"]*"' | cut -d'"' -f4)
    FOUND_LAT=$(echo "$NAME_BODY" | grep -o '"latitude":[0-9.-]*' | cut -d':' -f2)
    FOUND_LNG=$(echo "$NAME_BODY" | grep -o '"longitude":[0-9.-]*' | cut -d':' -f2)
    
    echo "   Found: $FOUND_NAME"
    echo "   City: $FOUND_CITY"
    echo "   Country: $FOUND_COUNTRY"
    echo "   Coordinates: $FOUND_LAT, $FOUND_LNG"
else
    echo "❌ Location from name failed (Status: $HTTP_STATUS)"
    echo "   Response: $NAME_BODY"
    exit 1
fi

echo ""

# Test 3: Test multiple location names
echo "3️⃣ Testing multiple location names..."
LOCATION_NAMES="${TEST_LOCATION_NAMES:-Amsterdam,London,Paris,Tokyo,New York}"

IFS=',' read -ra NAMES <<< "$LOCATION_NAMES"
for name in "${NAMES[@]}"; do
    name=$(echo "$name" | xargs) # trim whitespace
    echo "   Testing: $name"
    
    MULTI_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$API_BASE_URL/api/v1/location/name?name=$name")
    HTTP_STATUS=$(echo $MULTI_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    
    if [ "$HTTP_STATUS" -eq 200 ]; then
        FOUND_CITY=$(echo "$MULTI_RESPONSE" | grep -o '"city":"[^"]*"' | cut -d'"' -f4)
        FOUND_COUNTRY=$(echo "$MULTI_RESPONSE" | grep -o '"country":"[^"]*"' | cut -d'"' -f4)
        echo "     ✅ $name -> $FOUND_CITY, $FOUND_COUNTRY"
    else
        echo "     ❌ $name -> Failed (Status: $HTTP_STATUS)"
    fi
done

echo ""
echo "🎉 Location test completed successfully!"
echo "Location services are working correctly."
