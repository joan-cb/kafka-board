package helpers

import (
	"os"
)

// GetServerAddress returns the server address with port from environment
func GetServerAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9080"
	}
	return ":" + port
}
