package main

import (
	"context"
	"kafka-board/confluentRegistryAPI"
	"kafka-board/handlers"
	"kafka-board/helpers"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Logger available to all files in the package
var logger *slog.Logger

func main() {
	// Initialize logger
	logger = helpers.SetupLogger()

	// Create server with timeouts
	server := &http.Server{
		Addr:         helpers.GetServerAddress(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Initialize handler with logger
	handler := handlers.ReturnHandler(logger, confluentRegistryAPI.ReturnRegistryAPI(logger))

	// Set up routes
	http.HandleFunc("/", handler.HandleHomePage)
	http.HandleFunc("/schema/", handler.HandleSchemaPage)
	http.HandleFunc("/test-schema/", handler.HandleTestSchema)
	http.HandleFunc("/health", handler.HandleHealthCheck)
	http.HandleFunc("/test-payload", handler.HandleValidatePayload)

	// Channel to listen for errors coming from the listener.
	serverErrors := make(chan error, 1)

	// Start the server
	go func() {
		logger.Info("Server starting", "address", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Channel to listen for an interrupt or terminate signal from the OS.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		logger.Error("Server error",
			"error", err)

	case sig := <-shutdown:
		logger.Info("Server is shutting down",
			"signal", sig.String())

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Asking listener to shut down and shed load.
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Could not stop server gracefully", "error", err)
			server.Close()
		}
	}
}
