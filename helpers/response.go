package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

// CreatePayloadResponse creates a payloadTestResponse struct instance
// func CreatePayloadResponse(isValid bool, message string, statusCode int) types.PayloadTestResponse {
// 	return types.PayloadTestResponse{
// 		IsCompatible: isValid,
// 		Message:      message,
// 		StatusCode:   statusCode,
// 	}
// }

// CreateSchemaRegistryResponse creates a schemaRegistryResponse struct instance
func CreateResponse(isCompatible *bool, message string, httpStatus int, errorCode int) types.Response {
	return types.Response{
		IsCompatible: isCompatible,
		Message:      message,
		StatusCode:   httpStatus,
		ErrorCode:    errorCode,
	}
}

// Define specific response types
type CompatibilityResponse struct {
	IsCompatible bool `json:"is_compatible"`
}

type RegistryErrorResponse struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"message"`
}

// Union type for all possible responses
type RegistryResponse struct {
	Compatibility *CompatibilityResponse // For success
	Error         *RegistryErrorResponse // For error
	RawBody       []byte                 // Raw response
	StatusCode    int                    // HTTP status
}

// Process response based on status code
func ProcessResponse(resp *http.Response) (RegistryResponse, error) {
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		return RegistryResponse{
			StatusCode: resp.StatusCode,
		}, err
	}

	result := RegistryResponse{
		RawBody:    body,
		StatusCode: resp.StatusCode,
	}

	// Success case
	if resp.StatusCode == http.StatusOK {
		var compatibility CompatibilityResponse
		if err := json.Unmarshal(body, &compatibility); err != nil {
			log.Println(err)
			return result, err
		}
		result.Compatibility = &compatibility
		return result, nil
	}

	// Error case
	var errorResp RegistryErrorResponse
	if err := json.Unmarshal(body, &errorResp); err != nil {
		log.Println(err)
		return result, nil
	}

	result.Error = &errorResp
	return result, nil
}

// Convert to application response
func (r RegistryResponse) ToSchemaRegistryResponse() types.Response {
	// Success case
	if r.Compatibility != nil {
		isCompatible := r.Compatibility.IsCompatible
		return types.Response{
			IsCompatible: &isCompatible,
			Message:      "Schema is compatible",
			StatusCode:   r.StatusCode,
			ErrorCode:    0,
		}
	}

	// Error case
	if r.Error != nil {
		return types.Response{
			IsCompatible: nil,
			Message:      r.Error.Message,
			StatusCode:   r.StatusCode,
			ErrorCode:    r.Error.ErrorCode,
		}
	}

	// Unknown case
	return types.Response{
		IsCompatible: nil,
		Message:      fmt.Sprintf("Unexpected response (status: %d)", r.StatusCode),
		StatusCode:   r.StatusCode,
		ErrorCode:    r.StatusCode,
	}
}
