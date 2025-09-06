#!/bin/sh

# Check if binaries exist, if not build them
if [ ! -f "./bin/grpc-server" ] || [ ! -f "./bin/http-gateway" ]; then
    echo "Binaries not found. Building..."
    make build
fi

# Start gRPC server in background
./bin/grpc-server &
GRPC_PID=$!

# Wait for gRPC server to be ready
echo "Waiting for gRPC server to start..."
for i in $(seq 1 30); do
    if nc -z localhost 50051 2>/dev/null; then
        echo "gRPC server is ready!"
        break
    fi
    if [ $i -eq 30 ]; then
        echo "gRPC server failed to start within 30 seconds"
        exit 1
    fi
    sleep 1
done

# Start HTTP gateway in foreground
./bin/http-gateway &
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
