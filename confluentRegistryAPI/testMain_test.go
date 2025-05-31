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

	flag.Parse()

	// Check if integration tests should be skipped
	if testing.Short() || *skipServerTests {
		fmt.Println("Skipping integration tests that require a Schema Registry server")
		os.Exit(m.Run()) // Run the tests that don't require a server connection
	}

	// 1. Setup
	// registryAPI := ReturnRegistryAPI(testLogger)

	// // 2. Create test subjects
	// for _, testSubject := range newSubjects {
	// 	err := registryAPI.createTestSubject(testSubject)

	// 	if err != nil {
	// 		slog.Error("TestMain - Error creating test subject",
	// 			"error", err)
	// 		os.Exit(1)
	// 	}
	// }
	// // 3. Add config
	// for _, testSubject := range newSubjects {
	// 	err := registryAPI.createConfig(testSubject)
	// 	if err != nil {
	// 		slog.Error("Error creating config",
	// 			"error", err)
	// 		os.Exit(1)
	// 	}
	// }
	// 4. Run tests
	code := m.Run()
	slog.Debug("TestMain - Test results",
		"code", code)

	// 5. Teardown
	// cleanupTestEnvironment()

	// 6. Exit with the test result code
	os.Exit(code)
}
