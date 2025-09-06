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
CLOUDRUN_BINARY = bin/cloudrun-service

# Default target
.PHONY: all
all: clean deps proto build

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing Go dependencies..."
	$(GOGET) google.golang.org/grpc@latest
	$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GOGET) google.golang.org/api/drive/v3@latest
	$(GOGET) golang.org/x/oauth2/google@latest
	$(GOGET) google.golang.org/api/option@latest
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies installed successfully!"

# Install protobuf compiler and Go plugins
.PHONY: install-proto-deps
install-proto-deps:
	@echo "Installing protobuf dependencies..."
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

# Build cloudrun service
.PHONY: build
build: build-cloudrun

# Build cloudrun service
.PHONY: build-cloudrun
build-cloudrun:
	@echo "Building cloudrun service..."
	@mkdir -p bin
	$(GOBUILD) -o $(CLOUDRUN_BINARY) ./cmd/cloudrun

# Run cloudrun service
.PHONY: run
run: build-cloudrun
	./$(CLOUDRUN_BINARY)

# Run cloudrun service in background
.PHONY: run-bg
run-bg: build-cloudrun
	@echo "Starting cloudrun service in background..."
	./$(CLOUDRUN_BINARY) &

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
	@rm -f $(CLOUDRUN_BINARY)
	@rm -rf bin/
	@rm -f coverage.out coverage.html

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Initialize project (run once)
.PHONY: init
init: install-proto-deps deps proto
	@echo "Creating directory structure..."
	@mkdir -p cmd/cloudrun
	@mkdir -p internal/handlers internal/services
	@mkdir -p bin
	@echo "Project initialized successfully!"

# Development setup
.PHONY: dev-setup
dev-setup: init
	@echo "Setting up development environment..."
	@echo "Please add your Google credentials to credentials.json"

# Clean everything and rebuild
.PHONY: rebuild
rebuild: clean deps proto build

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make all              - Clean, install deps, generate proto, and build everything"
	@echo "  make init             - Initialize project structure (run once)"
	@echo "  make deps             - Install Go dependencies"
	@echo "  make proto            - Generate protobuf files"
	@echo "  make build            - Build cloudrun service"
	@echo "  make build-cloudrun   - Build cloudrun service only"
	@echo "  make run              - Run cloudrun service"
	@echo "  make run-bg           - Run cloudrun service in background"
	@echo "  make test             - Run tests"
	@echo "  make test-coverage    - Run tests with coverage"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make rebuild          - Clean everything and rebuild"
	@echo "  make fmt              - Format code"
	@echo "  make lint             - Lint code"
	@echo "  make help             - Show this help"