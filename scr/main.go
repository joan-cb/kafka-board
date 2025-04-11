package main

import (
	"log"
	"net/http"
)

func main() {
	// Set up HTTP handlers
	http.HandleFunc("/", handleHomePage)
	http.HandleFunc("/schema/", handleSchemaPage)
	http.HandleFunc("/test-schema/", handleTestSchema)

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Health check received")
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/test-payload", handleValidatePayload)

	// Start server
	log.Println("Server starting on http://localhost:9080")
	if err := http.ListenAndServe(":9080", nil); err != nil {
		log.Fatal(err)
	}
}
