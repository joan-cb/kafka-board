package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

// createTestSchemaRequest creates a new HTTP request for testing schema compatibility
func createTestSchemaRequest(subjectName string, version int, testJSON string) (*http.Request, error) {
	requestURL := baseRegistryURL + "/compatibility/subjects/" + subjectName + "/versions/" + strconv.Itoa(version)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(testJSON)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// makeHTTPRequest executes an HTTP request and returns the response
func makeHTTPRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return resp, nil
}

// readResponseBody reads and returns the response body
func readResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return body, nil
}

// processResponse handles the response based on its status code
func processCompatibilityResponse(body []byte, statusCode int) (schemaRegistryResponse, error) {
	switch statusCode {
	case http.StatusNotFound, http.StatusUnprocessableEntity:
		return handleErrorResponse(body, statusCode)
	case http.StatusInternalServerError:
		// For internal server errors, set IsCompatible to nil since we couldn't determine compatibility
		return schemaRegistryResponse{
			IsCompatible: nil,
			ErrorCode:    statusCode,
			Message:      "internal server error - please try again later",
			StatusCode:   statusCode,
		}, nil
	case http.StatusOK:
		return handleSuccessResponse(body, statusCode)
	default:
		// For unexpected status codes, set IsCompatible to nil since we couldn't determine compatibility
		return schemaRegistryResponse{
			IsCompatible: nil,
			ErrorCode:    statusCode,
			Message:      fmt.Sprintf("unexpected status code: %d, response: %s", statusCode, string(body)),
			StatusCode:   statusCode,
		}, nil
	}
}

// handleErrorResponse processes error responses from the schema registry
func handleErrorResponse(body []byte, statusCode int) (schemaRegistryResponse, error) {
	var errorResponse struct {
		ErrorCode int    `json:"error_code"`
		Message   string `json:"message"`
	}
	if err := json.Unmarshal(body, &errorResponse); err != nil {
		log.Printf("Error parsing error response: %v, body: %s", err, string(body))
		resp := createSchemaRegistryResponse(nil, fmt.Sprintf("error parsing error response: %v", err), statusCode, statusCode)
		return resp, nil
	}

	// Ensure message has a non-empty string value
	if errorResponse.Message == "" {
		errorResponse.Message = "None"
	}
	// For expected error responses, we set IsCompatible to false
	falseVal := false
	resp := createSchemaRegistryResponse(
		&falseVal,
		errorResponse.Message,
		statusCode,
		errorResponse.ErrorCode,
	)

	return resp, nil
}

// handleSuccessResponse processes successful responses from the schema registry
func handleSuccessResponse(body []byte, statusCode int) (schemaRegistryResponse, error) {
	var result struct {
		IsCompatible bool `json:"is_compatible"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		resp := createSchemaRegistryResponse(nil, fmt.Sprintf("error parsing response: %v", err), statusCode, statusCode)
		return resp, nil
	}

	// Return a pointer to the boolean value with default message
	resp := createSchemaRegistryResponse(&result.IsCompatible, "None", statusCode, 0)

	return resp, nil
}

// helper function to transform JSON to Schema Registry format
func transformJSONToSchemaFormat(jsonStr string) (string, error) {
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
		return "", fmt.Errorf("invalid JSON input: %v", err)
	}

	// Create the schema registry format
	schemaRegistryFormat := struct {
		Schema     string `json:"schema"`
		SchemaType string `json:"schemaType"`
	}{
		Schema:     jsonStr,
		SchemaType: "JSON",
	}

	// Marshal back to JSON
	result, err := json.Marshal(schemaRegistryFormat)
	if err != nil {
		return "", fmt.Errorf("error formatting schema: %v", err)
	}

	return string(result), nil
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, payload any) {
	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Set the status code
	w.WriteHeader(statusCode)

	// Encode and send the response
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// If encoding fails, send a simple error
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// createPayloadResponse creates a payloadTestResponse struct instance with the given parameters
func createPayloadResponse(isValid bool, message string, statusCode int) payloadTestResponse {
	return payloadTestResponse{
		IsCompatible: isValid,
		Message:      message,
		StatusCode:   statusCode,
	}
}

// createSchemaRegistryResponse creates a schemaRegistryResponse struct instance with the given parameters
func createSchemaRegistryResponse(isCompatible *bool, message string, httpStatus int, errorCode int) schemaRegistryResponse {
	return schemaRegistryResponse{
		IsCompatible: isCompatible,
		Message:      message,
		StatusCode:   httpStatus,
		ErrorCode:    errorCode,
	}
}

// checkErr is a helper function to check if an error is present
func checkErr(e error) bool {
	return e != nil
}
