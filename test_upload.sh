#!/bin/bash

# Test script for upload endpoint
echo "Testing upload endpoint..."

# Create a simple test image file
echo "Creating test image..."
convert -size 100x100 xc:red test_image.jpg 2>/dev/null || echo "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChwGA60e6kgAAAABJRU5ErkJggg==" | base64 -d > test_image.jpg

# Test the upload endpoint
echo "Testing POST /api/v1/images/upload..."
curl -X POST \
  -F "image=@test_image.jpg" \
  -F "title=Test Image" \
  -F "description=Test upload" \
  -F "latitude=37.7749" \
  -F "longitude=-122.4194" \
  -F "location_name=San Francisco" \
  -F "country=USA" \
  -F "city=San Francisco" \
  http://localhost:8080/api/v1/images/upload

echo -e "\n\nTest completed. Cleaning up..."
rm -f test_image.jpg
