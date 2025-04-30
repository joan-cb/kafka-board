#!/bin/bash

# Script to run tests for kafka-board project

# Function to check if a service is running on a specific port
check_service() {
  local host=$1
  local port=$2
  nc -z $host $port >/dev/null 2>&1
  return $?
}

# Set default values
USE_DOCKER=${USE_DOCKER:-false}
SCHEMA_REGISTRY_HOST=${SCHEMA_REGISTRY_HOST:-localhost}
SCHEMA_REGISTRY_PORT=${SCHEMA_REGISTRY_PORT:-8081}

# Parse command line arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    --docker)
      USE_DOCKER=true
      shift
      ;;
    --host=*)
      SCHEMA_REGISTRY_HOST="${1#*=}"
      shift
      ;;
    --port=*)
      SCHEMA_REGISTRY_PORT="${1#*=}"
      shift
      ;;
    --help)
      echo "Usage: $0 [options]"
      echo "Options:"
      echo "  --docker             Use Docker Compose for testing"
      echo "  --host=HOST          Schema Registry host (default: localhost)"
      echo "  --port=PORT          Schema Registry port (default: 8081)"
      echo "  --help               Show this help message"
      exit 0
      ;;
    *)
      echo "Unknown option: $1"
      echo "Run '$0 --help' for usage information"
      exit 1
      ;;
  esac
done

# Use Docker if requested
if [ "$USE_DOCKER" = "true" ]; then
  echo "Running tests with Docker Compose..."
  make test-docker
  exit $?
fi

# Check if Schema Registry is running
echo "Checking if Schema Registry is available at $SCHEMA_REGISTRY_HOST:$SCHEMA_REGISTRY_PORT..."
if check_service $SCHEMA_REGISTRY_HOST $SCHEMA_REGISTRY_PORT; then
  echo "Schema Registry is running at $SCHEMA_REGISTRY_HOST:$SCHEMA_REGISTRY_PORT"
  
  # Set environment variable for tests
  export REGISTRY_BASE_URL="http://$SCHEMA_REGISTRY_HOST:$SCHEMA_REGISTRY_PORT"
  
  # Run tests
  echo "Running tests with Schema Registry at $REGISTRY_BASE_URL..."
  go test ./...
  exit $?
else
  echo "Schema Registry is not available at $SCHEMA_REGISTRY_HOST:$SCHEMA_REGISTRY_PORT"
  echo ""
  echo "You have several options to run tests:"
  echo "1. Start Schema Registry locally on port 8081"
  echo "2. Use Docker for testing with: ./run-tests.sh --docker"
  echo "3. Specify a different host/port with: ./run-tests.sh --host=<host> --port=<port>"
  echo "4. Run tests in short mode (skipping integration tests) with: go test -short ./..."
  echo ""
  
  read -p "Do you want to run tests in short mode? (y/n) " answer
  if [[ "$answer" =~ ^[Yy]$ ]]; then
    echo "Running tests with -short flag to skip integration tests..."
    go test -short ./...
    exit $?
  else
    echo "Tests not run. Please set up Schema Registry and try again."
    exit 1
  fi
fi 