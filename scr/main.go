package main

import (
	"log"
	"net/http"
)

func main() {
	// Set up HTTP handlers

	handler := returnHandler(&registryAPI{})
	// No options for now, can be extended later
	http.HandleFunc("/", handler.handleHomePage)
	http.HandleFunc("/schema/", handler.handleSchemaPage)
	http.HandleFunc("/test-schema/", handler.handleTestSchema)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check received")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/test-payload", handler.handleValidatePayload)

	// Start server
	log.Println("Server starting on http://localhost:9080")
	if err := http.ListenAndServe(":9080", nil); err != nil {
		log.Fatal(err)
	}
}
