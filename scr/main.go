package main

import (
	"embed"
	"encoding/json"
	"log"
	"net/http"
)

//go:embed images/header.png images/slack-channel.png images/footer.png images/back.png images/github.png
var staticFiles embed.FS

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
	http.HandleFunc("/validate-payload", handleValidatePayload)
	http.HandleFunc("/validate", handleValidation)

	// Single handler for all static images
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		// Extract filename from URL path
		filename := r.URL.Path[len("/static/"):]

		// Read the embedded file (prepend with images/ directory)
		imageData, err := staticFiles.ReadFile("images/" + filename)
		if err != nil {
			log.Printf("Error reading static file %s: %v", filename, err)
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Set common headers
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "public, max-age=31536000")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Write the file content
		w.Write(imageData)
	})

	// Start server
	log.Println("Server starting on http://localhost:9080")
	if err := http.ListenAndServe(":9080", nil); err != nil {
		log.Fatal(err)
	}
}

func handleValidation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Schema  string `json:"schema"`
		Payload string `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Perform your validation logic here
	// This is where you'd validate the payload against the schema

	// Send back the validation result
	response := map[string]interface{}{
		"isValid": true,                    // Replace with actual validation result
		"message": "Validation successful", // Replace with actual message
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
