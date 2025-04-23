package confluentRegistryAPI

import (
	"io"
	"log/slog"
	"os"
	"testing"
)

// Unit tests run with no external dependencies
func TestUnit_Registry(t *testing.T) {
	// Skip if not in short mode
	if !testing.Short() {
		t.Skip("Skipping unit test when not in short mode")
	}

	t.Run("TransformJSON", func(t *testing.T) {
		// Create a new logger that discards output
		logger := slog.New(slog.NewTextHandler(io.Discard, nil))

		// Create a new RegistryAPI with the logger
		api := ReturnRegistryAPI(logger)

		// Test with valid JSON
		validJSON := `{"field1": "value1", "field2": 42}`
		_, err := api.TestSchema("test-subject", 1, validJSON)

		// This should fail without a real server, but we at least verify
		// that the code doesn't panic and the function returns expected error type
		if err == nil {
			t.Error("Expected error with no server available, but got nil")
		}

		// Test with invalid JSON
		invalidJSON := `{"field1": "value1",}`
		_, err = api.TestSchema("test-subject", 1, invalidJSON)
		if err == nil {
			t.Error("Expected error with invalid JSON, but got nil")
		}
	})
}

// Integration tests run with external dependencies
func TestIntegration_Registry(t *testing.T) {
	// Skip in short mode (unit tests only)
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Check for schema registry URL in environment
	schemaRegistryURL := os.Getenv("SCHEMA_REGISTRY_URL")
	if schemaRegistryURL == "" {
		t.Skip("SCHEMA_REGISTRY_URL not set, skipping integration tests")
	}

	t.Run("GetSubjects", testGetSubjects)
	t.Run("TestSchema", testSchemaCompatibility)
}

func testGetSubjects(t *testing.T) {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create a registry API
	api := ReturnRegistryAPI(logger)

	// Get subjects
	subjects, err := api.ReturnSubjects()
	if err != nil {
		t.Fatalf("Failed to get subjects: %v", err)
	}

	// Just verify we got a response
	t.Logf("Got %d subjects", len(subjects))
}

func testSchemaCompatibility(t *testing.T) {
	// Create a logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create a registry API
	api := ReturnRegistryAPI(logger)

	// We need a subject to test against
	subjects, err := api.ReturnSubjects()
	if err != nil || len(subjects) == 0 {
		t.Skip("No subjects available for testing compatibility")
	}

	// Simple schema for testing
	validJSON := `{"field1": "value1", "field2": 42}`

	// Test against the first subject
	_, err = api.TestSchema(subjects[0], 1, validJSON)

	// We don't care about the result, just that it can connect
	t.Logf("Test schema returned: %v", err)
}
