package helpers

import (
	"log/slog"
	"os"
)

// SetupLogger initializes and returns a structured logger
func SetupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: GetLogLevel(),
	}))
}

// GetLogLevel returns the appropriate log level based on environment variable
func GetLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
