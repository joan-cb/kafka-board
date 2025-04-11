package main

import (
	"log"
	"net/http"
)

func main() {
	// Set up HTTP handlers
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/schema/", handleSchema)
	http.HandleFunc("/test-schema/", handleTestSchema)
	http.HandleFunc("/test-schema", handleTestSchemaAPI)
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
