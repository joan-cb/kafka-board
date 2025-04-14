package main

import (
	"kafka-board/handlers"
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

	handler := handlers.ReturnHandler(logger)
	// No options for now, can be extended later
	http.HandleFunc("/", handler.HandleHomePage)
	http.HandleFunc("/schema/", handler.HandleSchemaPage)
	http.HandleFunc("/test-schema/", handler.HandleTestSchema)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Health check received")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/test-payload", handler.HandleValidatePayload)

	// Start server
	logger.Info("Server starting on port 9080")
	if err := http.ListenAndServe(":9080", nil); err != nil {
		log.Fatal(err)
	}
}
