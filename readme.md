# Kafka Schema Registry UI

A web-based interface for managing and validating schemas in a Confluent Schema Registry.

## Overview

This application provides a user-friendly web interface to interact with a Confluent Schema Registry. It allows users to:
- View all registered subjects (schemas)
- Inspect schema details and versions
- Test schema compatibility
- Validate JSON payloads against schemas
- Retrieve global and subject-level configuration

## Architecture

The application is built with Go and consists of three main components:

1. **Main Application (`main.go`)**
   - Sets up HTTP server and routes
   - Initializes logging
   - Handles health checks
   - Serves on port 9080

2. **Handlers (`handlers/handlers.go`)**
   - Manages HTTP endpoints and request handling
   - Implements UI rendering with Go templates
   - Handles schema testing and payload validation
   - Provides error handling and logging

3. **Registry API Client (`confluentRegistryAPI/apiCalls.go`)**
   - Communicates with Confluent Schema Registry
   - Implements REST API calls to the registry
   - Handles JSON transformations and validations
   - Manages schema compatibility testing

## Key Features

### Schema Management
- View all registered subjects
- Browse schema versions
- View schema details and configurations
- Pretty-print JSON schemas

### Schema Testing
- Test schema compatibility
- Validate JSON payloads against schemas
- Get detailed validation error messages
- Support for different compatibility modes

### Configuration
- View and manage global configuration
- Subject-level configuration management
- Compatibility mode settings
- Default configuration handling

## Technical Details

- Built with Go 1.22+
- Uses standard library `net/http` for web server
- Implements structured logging with `slog`
- JSON schema validation with `gojsonschema`
- Template-based UI rendering
- REST API communication with Schema Registry

## Dependencies

- Go 1.22 or later
- Confluent Schema Registry (v7.3.0+)
- Docker (for containerized deployment)

## API Endpoints

- `/` - Home page with subject list
- `/schema/` - Schema details page
- `/test-schema/` - Schema testing interface
- `/test-payload` - Payload validation endpoint
- `/health` - Health check endpoint

## Configuration

The application connects to the Schema Registry at `http://schema-registry:8081` by default. This can be configured through environment variables.

## Logging

The application uses structured logging with different levels:
- Debug: Detailed information for debugging
- Info: General operational information
- Error: Error conditions that need attention

## Error Handling

- Comprehensive error handling throughout the application
- User-friendly error messages
- Detailed logging for debugging
- Proper HTTP status codes

## Security

- Input validation for all user inputs
- Safe JSON parsing
- Error handling for malformed requests
- Proper HTTP status codes for different scenarios
