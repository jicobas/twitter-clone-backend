.PHONY: run-dev test test-coverage build clean docker-build docker-run

# Variables
APP_NAME=twitter-clone-backend
DOCKER_IMAGE=twitter-clone-backend

# Development
run-dev:
	@echo "Running in development mode (in-memory storage)..."
	go run cmd/server/main.go

# Testing
test:
	@echo "Running tests..."
	go test ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -cover ./...

test-integration:
	@echo "Running integration tests..."
	go test -tags=integration ./...

# Build
build:
	@echo "Building application..."
	go build -o bin/$(APP_NAME) cmd/server/main.go

clean:
	@echo "Cleaning build files..."
	rm -rf bin/

# Docker
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run:
	@echo "Running with Docker Compose..."
	docker-compose up --build

# Dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy

# Format and lint
fmt:
	@echo "Formatting code..."
	go fmt ./...

vet:
	@echo "Running go vet..."
	go vet ./...
