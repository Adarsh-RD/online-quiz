.PHONY: all build run test lint clean docker-build docker-up docker-down

APP_NAME=online-quiz
BIN_DIR=bin

all: lint test build

build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/server/main.go
	@echo "Build complete."

run:
	@echo "Running $(APP_NAME)..."
	@go run ./cmd/server/main.go

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)
	@go clean

docker-build:
	@echo "Building Docker image..."
	@docker-compose build

docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d

docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
