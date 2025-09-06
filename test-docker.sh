#!/bin/bash

echo "🐳 Testing Docker container locally"
echo ""

# Build the Docker image
echo "📦 Building Docker image..."
docker build -t portfolio-images-api .

# Run the container
echo "🚀 Starting container..."
docker run -d \
  --name portfolio-images-test \
  -p 8080:8080 \
  -p 50051:50051 \
  -e GOOGLE_DRIVE_FOLDER_ID="$GOOGLE_DRIVE_FOLDER_ID" \
  -e GOOGLE_MAPS_API_KEY="$GOOGLE_MAPS_API_KEY" \
  -e CLOUD_SQL_CONNECTION_NAME="$CLOUD_SQL_CONNECTION_NAME" \
  -e CLOUD_SQL_DATABASE="$CLOUD_SQL_DATABASE" \
  -e CLOUD_SQL_USER="$CLOUD_SQL_USER" \
  -e CLOUD_SQL_PASSWORD="$CLOUD_SQL_PASSWORD" \
  -e GRPC_PORT="50051" \
  portfolio-images-api

# Wait for services to start
echo "⏳ Waiting for services to start..."
sleep 10

# Test health endpoint
echo "🔍 Testing health endpoint..."
curl -s http://localhost:8080/health | jq .

# Test image count
echo "🔍 Testing image count endpoint..."
curl -s http://localhost:8080/api/v1/images/count | jq .

echo ""
echo "✅ Container is running!"
echo "🌐 API available at: http://localhost:8080"
echo "🔧 gRPC server on: localhost:50051"
echo ""
echo "To stop the container:"
echo "docker stop portfolio-images-test"
echo "docker rm portfolio-images-test"
