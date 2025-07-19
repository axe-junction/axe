# Gateway Service Makefile

.PHONY: help build run dev test clean docker-build docker-run docker-dev proto-gen

# Default target
help:
	@echo "Available commands:"
	@echo "  build       - Build the gateway binary"
	@echo "  run         - Run the gateway service"
	@echo "  dev         - Run with hot reload (requires air)"
	@echo "  test        - Run tests"
	@echo "  clean       - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run  - Run with Docker Compose"
	@echo "  docker-dev  - Run with Docker Compose (dev profile)"
	@echo "  proto-gen   - Generate protobuf code"
	@echo "  logs        - Show Docker logs"
	@echo "  down        - Stop Docker services"

# Build the binary
build:
	@echo "🔨 Building gateway service..."
	go build -o bin/gateway ./cmd/main.go

# Run the service locally
run: build
	@echo "🚀 Running gateway service..."
	./bin/gateway

# Development with hot reload (requires air: go install github.com/air-verse/air@latest)
dev:
	@echo "🔥 Running gateway service with hot reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "Air not found. Install with: go install github.com/air-verse/air@latest"; \
		echo "Or run: make run"; \
	fi

# Run tests
test:
	@echo "🧪 Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	rm -rf bin/
	go clean

# Generate protobuf code
proto-gen:
	@echo "📦 Generating protobuf code..."
	protoc --go_out=. --go-grpc_out=. pkg/pb/*.proto

# Docker commands
docker-build:
	@echo "🐳 Building Docker image..."
	docker-compose build

docker-run:
	@echo "🐳 Running with Docker Compose..."
	docker-compose up -d

docker-dev:
	@echo "🐳 Running with Docker Compose (dev mode)..."
	docker-compose --profile dev up -d

logs:
	@echo "📋 Showing Docker logs..."
	docker-compose logs -f

down:
	@echo "⬇️ Stopping Docker services..."
	docker-compose down

# Database commands
db-reset:
	@echo "🔄 Resetting database..."
	docker-compose down -v
	docker-compose up -d gateway-db
	@echo "⏳ Waiting for database to be ready..."
	sleep 10
	docker-compose up -d gateway

# Check service health
health:
	@echo "🏥 Checking service health..."
	@curl -s http://localhost:8080/health | jq . || echo "Service not running or jq not installed"

# Install dependencies
deps:
	@echo "📦 Installing dependencies..."
	go mod download
	go mod tidy

# Format code
fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...
	goimports -w .

# Lint code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi
