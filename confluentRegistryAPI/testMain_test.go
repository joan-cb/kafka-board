package confluentRegistryAPI

// import (
// 	"log/slog"
// 	"net/http"
// 	"os"
// 	"testing"
// 	"time"
// )

// // TestMain controls the entire test execution
// func TestMain(m *testing.M) {
// 	// Only set up integration test resources if not in short mode
// 	if !testing.Short() {
// 		setupIntegrationTests()
// 		defer teardownIntegrationTests()
// 	}

// 	// Run all tests and exit
// 	os.Exit(m.Run())
// }

// // setupIntegrationTests prepares the integration test environment
// func setupIntegrationTests() {
// 	// Read server URL from environment variable, with fallback
// 	serverURL := os.Getenv("SCHEMA_REGISTRY_URL")
// 	if serverURL == "" {
// 		serverURL = "http://schema-registry:8081" // Default for local development
// 	}

// 	// Override the global baseRegistryURL
// 	baseRegistryURL = serverURL

// 	// Initialize logger for tests
// 	logggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
// 		Level: slog.LevelDebug,
// 	})
// 	logger = slog.New(logggerHandler)

// 	// Wait for server to be available
// 	waitForServer(baseRegistryURL, 30*time.Second)

// 	// Optional: Set up test data if needed
// 	setupTestData()
// }

// // teardownIntegrationTests cleans up after integration tests
// func teardownIntegrationTests() {
// 	// Clean up test data if needed
// 	cleanupTestData()
// }

// // waitForServer attempts to connect to the server until it's available or timeout
// func waitForServer(url string, timeout time.Duration) {
// 	deadline := time.Now().Add(timeout)
// 	for time.Now().Before(deadline) {
// 		resp, err := http.Get(url + "/subjects")
// 		if err == nil {
// 			resp.Body.Close()
// 			if resp.StatusCode == http.StatusOK {
// 				return
// 			}
// 		}
// 		time.Sleep(1 * time.Second)
// 	}
// 	panic("Schema Registry server not available within timeout")
// }
