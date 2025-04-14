package confluentRegistryAPI

import (
	"log/slog"
	"testing"
)

func TestRegistryAPI_Integration(t *testing.T) {
	// Skip if running in short mode
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	t.Run("TestReturnSubjects", testReturnSubjects)
}

func testReturnSubjects(t *testing.T) {
	slog.Debug("testReturnSubjects")
}
