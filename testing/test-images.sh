#!/bin/bash

# Simple Images Test
# Tests upload endpoint and then gets the uploaded image

echo "📸 Images Upload & Get Test"
echo "============================"

# Load environment variables from test.env if it exists
if [ -f "testing/test.env" ]; then
    echo "📋 Loading test configuration from testing/test.env"
    export $(grep -v '^#' testing/test.env | xargs)
fi

# Configuration
API_BASE_URL="${API_BASE_URL:-https://background-image-drive-api-189526192204.us-west1.run.app}"
TEST_IMAGE="${TEST_IMAGE_PATH:-testing/test_image.jpg}"

echo "API Base URL: $API_BASE_URL"
echo "Test Image: $TEST_IMAGE"
echo ""

# Check if test image exists
if [ ! -f "$TEST_IMAGE" ]; then
    echo "❌ Test image not found at: $TEST_IMAGE"
    echo "   Please set TEST_IMAGE_PATH environment variable or place test image at testing/test_image.jpg"
    exit 1
fi

echo "✅ Test image found: $(ls -lh "$TEST_IMAGE" | awk '{print $5}')"
echo ""

# Test 1: Upload image
echo "1️⃣ Testing image upload..."
UPLOAD_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" -X POST \
  -F "image=@$TEST_IMAGE" \
  -F "title=${TEST_IMAGE_TITLE:-Test Image}" \
  -F "description=${TEST_IMAGE_DESCRIPTION:-Test upload}" \
  -F "latitude=${TEST_LATITUDE:-52.3676}" \
  -F "longitude=${TEST_LONGITUDE:-4.9041}" \
  -F "location_name=${TEST_LOCATION_NAME:-Amsterdam}" \
  -F "country=${TEST_COUNTRY:-Netherlands}" \
  -F "city=${TEST_CITY:-Amsterdam}" \
  "$API_BASE_URL/api/v1/images/upload")

HTTP_STATUS=$(echo $UPLOAD_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
UPLOAD_BODY=$(echo $UPLOAD_RESPONSE | sed -e 's/HTTPSTATUS:.*//g')

if [ "$HTTP_STATUS" -eq 201 ]; then
    echo "✅ Upload successful (Status: $HTTP_STATUS)"
    
    # Extract image ID from response (try both new and old formats)
    IMAGE_ID=$(echo "$UPLOAD_BODY" | grep -o '"image_id":"[^"]*"' | cut -d'"' -f4)
    if [ -z "$IMAGE_ID" ]; then
        # Fallback to old format (image ID in metadata)
        IMAGE_ID=$(echo "$UPLOAD_BODY" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
    fi
    
    if [ -n "$IMAGE_ID" ]; then
        echo "   Image ID: $IMAGE_ID"
    else
        echo "   ⚠️  Could not extract image ID from response"
        echo "   Response: $UPLOAD_BODY"
        exit 1
    fi
else
    echo "❌ Upload failed (Status: $HTTP_STATUS)"
    echo "   Response: $UPLOAD_BODY"
    exit 1
fi

echo ""

# Test 2: Get image by ID
if [ -n "$IMAGE_ID" ]; then
    echo "2️⃣ Testing get image by ID..."
    GET_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$API_BASE_URL/api/v1/images/$IMAGE_ID")
    
    HTTP_STATUS=$(echo $GET_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
    GET_BODY=$(echo $GET_RESPONSE | sed -e 's/HTTPSTATUS:.*//g')
    
    if [ "$HTTP_STATUS" -eq 200 ]; then
        echo "✅ Get image successful (Status: $HTTP_STATUS)"
        echo "   Image data retrieved successfully"
    else
        echo "❌ Get image failed (Status: $HTTP_STATUS)"
        echo "   Response: $GET_BODY"
        exit 1
    fi
else
    echo "⚠️  Skipping get image test - no image ID available"
fi

echo ""

# Test 3: List all images
echo "3️⃣ Testing list all images..."
LIST_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$API_BASE_URL/api/v1/images")
HTTP_STATUS=$(echo $LIST_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
LIST_BODY=$(echo $LIST_RESPONSE | sed -e 's/HTTPSTATUS:.*//g')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ List images successful (Status: $HTTP_STATUS)"
    IMAGE_COUNT=$(echo "$LIST_BODY" | grep -o '"id"' | wc -l)
    echo "   Found $IMAGE_COUNT images in database"
else
    echo "❌ List images failed (Status: $HTTP_STATUS)"
    echo "   Response: $LIST_BODY"
    exit 1
fi

echo ""

# Test 4: Get image count
echo "4️⃣ Testing image count..."
COUNT_RESPONSE=$(curl -s -w "HTTPSTATUS:%{http_code}" "$API_BASE_URL/api/v1/images/count")
HTTP_STATUS=$(echo $COUNT_RESPONSE | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
COUNT_BODY=$(echo $COUNT_RESPONSE | sed -e 's/HTTPSTATUS:.*//g')

if [ "$HTTP_STATUS" -eq 200 ]; then
    echo "✅ Image count successful (Status: $HTTP_STATUS)"
    echo "   Total images: $COUNT_BODY"
else
    echo "❌ Image count failed (Status: $HTTP_STATUS)"
    echo "   Response: $COUNT_BODY"
    exit 1
fi

echo ""
echo "🎉 Images test completed successfully!"
echo "Uploaded image ID: $IMAGE_ID"
