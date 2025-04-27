# Kafka Schema Registry UI

A web-based interface for managing and validating schemas in a Confluent Schema Registry.

## Overview

This application provides a user-friendly web interface to interact with a Confluent Schema Registry. It allows users to:
- View all registered subjects (schemas)
- Inspect schema details and versions
- Test schema compatibility
- Validate JSON payloads against schemas
- Retrieve global and subject-level configuration

## Implementation

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

## Technical Stack

- Built with Go 1.24
- Uses standard library `net/http` for web server
- Implements structured logging with `slog`
- JSON schema validation with `gojsonschema`
- REST API communication with Schema Registry
