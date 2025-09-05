# Variables
PROTO_DIR = proto
GENERATED_DIR = $(PROTO_DIR)/gen
GO_OUT = $(GENERATED_DIR)
PROTO_FILES = $(PROTO_DIR)/*.proto

# Go related variables
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod

# Binary names
SERVER_BINARY = bin/grpc-server
HTTP_BINARY = bin/http-gateway

# Default target
.PHONY: all
all: clean proto build

# Install dependencies
.PHONY: deps
deps:
	$(GOMOD) download
	$(GOMOD) tidy

# Install protobuf compiler and Go plugins
.PHONY: install-proto-deps
install-proto-deps:
	# Install protoc-gen-go and protoc-gen-go-grpc
	$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf files
.PHONY: proto
proto: clean-proto
	@echo "Generating protobuf files..."
	@mkdir -p $(GENERATED_DIR)
	protoc \
		--go_out=$(GO_OUT) \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) \
		--go-grpc_opt=paths=source_relative \
		--proto_path=$(PROTO_DIR) \
		$(PROTO_FILES)
	@echo "Protobuf files generated successfully!"

# Clean generated protobuf files
.PHONY: clean-proto
clean-proto:
	@echo "Cleaning generated protobuf files..."
	@rm -rf $(GENERATED_DIR)

# Build all binaries
.PHONY: build
build: build-server build-http

# Build gRPC server
.PHONY: build-server
build-server:
	@echo "Building gRPC server..."
	@mkdir -p bin
	$(GOBUILD) -o $(SERVER_BINARY) ./cmd/server

# Build HTTP gateway
.PHONY: build-http
build-http:
	@echo "Building HTTP gateway..."
	@mkdir -p bin
	$(GOBUILD) -o $(HTTP_BINARY) ./cmd/http

# Run gRPC server
.PHONY: run-server
run-server: build-server
	./$(SERVER_BINARY)

# Run HTTP gateway
.PHONY: run-http
run-http: build-http
	./$(HTTP_BINARY)

# Run both services (in background)
.PHONY: run-all
run-all:
	@echo "Starting gRPC server..."
	./$(SERVER_BINARY) &
	@echo "Starting HTTP gateway..."
	./$(HTTP_BINARY) &

# Test
.PHONY: test
test:
	$(GOTEST) -v ./...

# Test with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	@rm -f $(SERVER_BINARY)
	@rm -f $(HTTP_BINARY)
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Initialize project (run once)
.PHONY: init
init: install-proto-deps deps proto
	@echo "Creating directory structure..."
	@mkdir -p cmd/server cmd/http
	@mkdir -p internal/handlers internal/services
	@mkdir -p bin
	@echo "Project initialized successfully!"

# Development setup
.PHONY: dev-setup
dev-setup: init
	@echo "Setting up development environment..."
	@echo "Please add your Google credentials to credentials.json"

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make all              - Clean, generate proto, and build everything"
	@echo "  make init             - Initialize project structure (run once)"
	@echo "  make deps             - Install Go dependencies"
	@echo "  make proto            - Generate protobuf files"
	@echo "  make build            - Build all binaries"
	@echo "  make build-server     - Build gRPC server only"
	@echo "  make build-http       - Build HTTP gateway only"
	@echo "  make run-server       - Run gRPC server"
	@echo "  make run-http         - Run HTTP gateway"
	@echo "  make run-all          - Run both services"
	@echo "  make test             - Run tests"
	@echo "  make test-coverage    - Run tests with coverage"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make fmt              - Format code"
	@echo "  make help             - Show this help"