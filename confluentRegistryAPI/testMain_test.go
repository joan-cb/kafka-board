package confluentRegistryAPI

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"testing"
)

// skipServerTests flag controls whether to skip tests that require a real Schema Registry server
// Default is false, meaning we'll try to connect to a real server by default
var skipServerTests = flag.Bool("skip-server-tests", false, "Skip tests that require a real Schema Registry server")

//to do: create testStruct for test cases and not resuse registryAPI struct

func TestMain(m *testing.M) {
	// Set up logger with debug level
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	testLogger := slog.New(logHandler)
	slog.SetDefault(testLogger)

	// Set up environment for testing - use schema-registry:8081 as the default URL
	if os.Getenv("SCHEMA_REGISTRY_URL") == "" {
		os.Setenv("SCHEMA_REGISTRY_URL", "http://schema-registry:8081")
		testLogger.Info("Setting Schema Registry URL for tests",
			"url", os.Getenv("SCHEMA_REGISTRY_URL"))
	}

	flag.Parse()

	// Check if integration tests should be skipped
	if testing.Short() || *skipServerTests {
		fmt.Println("Skipping integration tests that require a Schema Registry server")
		os.Exit(m.Run()) // Run the tests that don't require a server connection
	}

	// 1. Setup
	registryAPI := ReturnRegistryAPI(testLogger)

	subjectNames, err := registryAPI.ReturnSubjects()

	if err != nil {
		slog.Error("Error getting subjects",
			"error", err)
		os.Exit(1)
	}
	slog.Debug("TestMain - Subjects returned",
		"subjects", subjectNames)
	// Delete all subjects
	deletedSubjects, err := registryAPI.deleteAllSubjects(subjectNames)
	slog.Debug("TestMain - Deleted subjects",
		"subjects", deletedSubjects)
	if err != nil {
		slog.Error("Error deleting test subject",
			"error", err)
		os.Exit(1)
	}
	// Create test subjects
	err = registryAPI.createTestSubject(newSubjects)
	if err != nil {
		slog.Error("Error creating test subject",
			"error", err)
		os.Exit(1)
	}

	// prepareTestEnvironment()

	// 2. Run tests
	code := m.Run()
	slog.Debug("TestMain - Test results",
		"code", code)

	// 3. Teardown
	// cleanupTestEnvironment()

	// 4. Exit with the test result code
	os.Exit(code)
}
