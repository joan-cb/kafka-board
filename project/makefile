# Variables
SRV_BINARY=kafkaBoard
GO_BUILD_FLAGS=GOOS=linux CGO_ENABLED=0

.PHONY: all build docker-up docker-down clean

## Build and run the service
all: build up

## Build the service
build:
	@echo "Building service..."
	cd ../scr && $(GO_BUILD_FLAGS) go build -o $(SRV_BINARY) .
	@echo "Build complete!"

## Stop Docker services
down:
	@echo "Stopping services..."
	docker compose down
	@echo "Services stopped!"

## Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f ../scr/$(SRV_BINARY)
	@echo "Clean complete!"

## Show help
help:
	@grep -E '^##' Makefile | sed 's/## //'

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker compose up --build -d
	@echo "Docker images started!"

## build_service: builds the service binary
build_service:
	@echo "Building service binary..."
	cd ../scr && env GOOS=linux CGO_ENABLED=0 go build -o $(SRV_BINARY) .
	@echo "Service binary built!"

## up_build: stops, builds, and starts containers
up_build: down build
	@echo "Building and starting Docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

rebuild-board:
	@echo "Stopping kafka-board container..."
	docker compose stop kafka-board
	@echo "Rebuilding kafka-board..."
	docker compose build kafka-board
	@echo "Starting kafka-board container..."
	docker compose up -d kafka-board
	@echo "Kafka-board rebuild complete!"

help:
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'