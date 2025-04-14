# Variables
SRV_BINARY=kafkaBoard
GO_BUILD_FLAGS=GOOS=linux CGO_ENABLED=0

.PHONY: all build docker-up docker-down clean

## Build and run the service
all: build up

## Build the service
build:
	@echo "Building service..."
	cd ../ && $(GO_BUILD_FLAGS) go build -o $(SRV_BINARY) .
	@echo "Build complete!"

## Stop Docker services
down:
	@echo "Stopping services..."
	docker compose down
	@echo "Services stopped!"

## Clean up build artifacts
clean:
	@echo "Cleaning up..."
	rm -f ../$(SRV_BINARY)
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

 