#!/bin/sh

# Start gRPC server in background
./grpc-server &
GRPC_PID=$!

# Wait a moment for gRPC server to start
sleep 2

# Start HTTP gateway in foreground
./http-gateway &
HTTP_PID=$!

# Function to handle shutdown
shutdown() {
    echo "Shutting down services..."
    kill $GRPC_PID $HTTP_PID 2>/dev/null
    wait $GRPC_PID $HTTP_PID
    exit 0
}

# Set up signal handlers
trap shutdown SIGTERM SIGINT

# Wait for either process to exit
wait $GRPC_PID $HTTP_PID
