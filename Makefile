# Variables
SRV_BINARY=kafkaBoard
GO_BUILD_FLAGS=GOOS=linux CGO_ENABLED=0

.PHONY: all build docker-up docker-down clean test test-docker

## Build and run the service
all: build up

## Build the service
build:
	@echo "Building service..."
	$(GO_BUILD_FLAGS) go build -o $(SRV_BINARY) .
	@echo "Build complete!"

## Stop Docker services
down:
	@echo "Stopping services..."
	docker compose down
	@echo "Services stopped!"

## Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f $(SRV_BINARY)
	@echo "Clean complete!"

## Show help
help:
	@grep -E '^##' Makefile | sed 's/## //'

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up --build -d
	@echo "Docker images started!"

## up_build: stops, builds, and starts containers
up_build: down build
	@echo "Building and starting Docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## Run tests in short mode
test-short:
	@echo "Running tests in short mode..."
	go test -short ./...
	@echo "Short tests complete!"

## Run all tests
test:
	@echo "Running all tests..."
	go test ./...
	@echo "All tests complete!"

## Run tests with Docker Compose setup for Schema Registry
test-docker:
	@echo "Running tests with Docker Compose..."
	docker compose -f docker-compose.test.yml up --build --abort-on-container-exit
	@echo "Docker tests complete!"
	docker compose -f docker-compose.test.yml down 