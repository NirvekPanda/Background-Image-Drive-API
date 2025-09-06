# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install git, ca-certificates, and SQLite development files
RUN apk add --no-cache git ca-certificates sqlite-dev gcc musl-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the Cloud Run application (CGO enabled for SQLite)
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o cloudrun-service ./cmd/cloudrun

# Final stage
FROM alpine:latest

# Install ca-certificates, netcat for health checks, and SQLite runtime
RUN apk --no-cache add ca-certificates netcat-openbsd sqlite

# Create non-root user
RUN adduser -D -s /bin/sh appuser

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/cloudrun-service .

# Create directories for OAuth2 credentials (will be mounted as secrets in production)
RUN mkdir -p /app/secrets

# Change ownership to non-root user
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Start the Cloud Run service
CMD ["./cloudrun-service"]
