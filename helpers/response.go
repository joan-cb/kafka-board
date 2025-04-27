package helpers

import (
	"encoding/json"
	"net/http"

	"kafka-board/types"
)

// SendJSONResponse sends a JSON response
func SendJSONResponse(w http.ResponseWriter, statusCode int, payload any) {
	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Set the status code
	w.WriteHeader(statusCode)

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// If encoding fails, send a simple error
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// CreateResponseObject builds a standardized response structure
// This function creates a consistent response format for HTTP handlers
// to use when responding to browser or API requests
func CreateResponseObject(isCompatible *bool, message string, httpStatus int, errorCode int) types.Response {
	return types.Response{
		IsCompatible: isCompatible,
		Message:      message,
		StatusCode:   httpStatus,
		ErrorCode:    errorCode,
	}
}
