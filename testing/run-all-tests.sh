#!/bin/bash

# Master Test Runner
# Runs all simple tests in sequence

echo "🧪 Portfolio Images API - Complete Test Suite"
echo "=============================================="
echo ""

# Load environment variables from test.env if it exists
if [ -f "testing/test.env" ]; then
    echo "📋 Loading test configuration from testing/test.env"
    export $(grep -v '^#' testing/test.env | xargs)
fi

# Configuration
API_BASE_URL="${API_BASE_URL:-https://background-image-drive-api-189526192204.us-west1.run.app}"

echo "API Base URL: $API_BASE_URL"
echo "Test Image: ${TEST_IMAGE_PATH:-testing/test_image.jpg}"
echo "Test Coordinates: ${TEST_LATITUDE:-not set}, ${TEST_LONGITUDE:-not set}"
echo "Test Locations: ${TEST_LOCATION_NAMES:-not set}"
echo ""

# Run tests in sequence
echo "🚀 Starting test sequence..."
echo ""

# Test 1: Health Check
echo "=== TEST 1: HEALTH CHECK ==="
./testing/test-health.sh
if [ $? -ne 0 ]; then
    echo ""
    echo "❌ Health check failed. Stopping tests."
    exit 1
fi

echo ""
echo "=== TEST 2: IMAGES TEST ==="
./testing/test-images.sh
if [ $? -ne 0 ]; then
    echo ""
    echo "❌ Images test failed. Stopping tests."
    exit 1
fi

echo ""
echo "=== TEST 3: LOCATION TEST ==="
./testing/test-location.sh
if [ $? -ne 0 ]; then
    echo ""
    echo "❌ Location test failed."
    exit 1
fi

echo ""
echo "🎉 ALL TESTS COMPLETED SUCCESSFULLY!"
echo "====================================="
echo "✅ Health check passed"
echo "✅ Images upload/get passed"
echo "✅ Location services passed"
echo ""
echo "API is fully functional and ready for production use."
