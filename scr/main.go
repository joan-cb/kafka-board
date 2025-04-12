package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
)

// Logger available to all files in the package
var logger *slog.Logger

func main() {
	// Set up HTTP handlers

	// Initialize logger before using it
	logggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug, // Only logs Info, Warn, and Error
	})
	logger = slog.New(logggerHandler)
	logger.Info("Starting server on port 9080")

	handler := returnHandler(&registryAPI{})
	// No options for now, can be extended later
	http.HandleFunc("/", handler.handleHomePage)
	http.HandleFunc("/schema/", handler.handleSchemaPage)
	http.HandleFunc("/test-schema/", handler.handleTestSchema)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health check received")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/test-payload", handler.handleValidatePayload)

	// Start server
	logger.Info("Server starting on port 9080")
	if err := http.ListenAndServe(":9080", nil); err != nil {
		log.Fatal(err)
	}
}
