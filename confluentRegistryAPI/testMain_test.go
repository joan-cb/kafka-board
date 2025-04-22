package confluentRegistryAPI

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"testing"
)

//to do: create testStruct for test cases and not resuse registryAPI struct

func TestMain(m *testing.M) {
	// Set up logger with debug level
	logHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	testLogger := slog.New(logHandler)
	slog.SetDefault(testLogger)

	// Set up environment for testing - use localhost when running tests
	if os.Getenv("SCHEMA_REGISTRY_URL") == "" {
		os.Setenv("SCHEMA_REGISTRY_URL", "http://localhost:8090")
		testLogger.Info("Setting Schema Registry URL for tests",
			"url", os.Getenv("SCHEMA_REGISTRY_URL"))
	}

	flag.Parse()
	if testing.Short() {
		fmt.Println("Skipping integration tests in short mode")
		os.Exit(0)
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
