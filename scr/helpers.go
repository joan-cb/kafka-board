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
	log.Printf("Creating test schema request for subject '%s', version %d", subjectName, version)

	requestURL := baseRegistryURL + "/compatibility/subjects/" + subjectName + "/versions/" + strconv.Itoa(version)
	log.Printf("Request URL: %s", requestURL)

	// Transform JSON to Schema Registry format
	payload, err := transformJSONToSchemaFormat(testJSON)
	if err != nil {
		log.Printf("Error transforming JSON to Schema Registry format: %v", err)
		return nil, fmt.Errorf("error transforming JSON: %v", err)
	}
	log.Printf("Transformed JSON returned by transformJSONToSchemaFormat: %s", payload)

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")
	req.Header.Set("Content-Type", "application/json")

	log.Printf("Request headers set: Accept=%s, Content-Type=%s",
		req.Header.Get("Accept"),
		req.Header.Get("Content-Type"))

	return req, nil
}

// makeHTTPRequest executes an HTTP request and returns the response
func makeHTTPRequest(req *http.Request) (*http.Response, error) {
	log.Printf("Making HTTP request to %s", req.URL.String())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error making HTTP request: %v", err)
		return nil, fmt.Errorf("error making request: %v", err)
	}

	log.Printf("Received response with status code: %d", resp.StatusCode)
	return resp, nil
}

// readResponseBody reads and returns the response body
func readResponseBody(resp *http.Response) ([]byte, error) {
	log.Printf("Reading response body for status code %d", resp.StatusCode)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	log.Printf("Successfully read response body of length %d bytes", len(body))
	return body, nil
}

// handleErrorResponse processes error responses from the schema registry
func handleErrorResponse(body []byte, statusCode int) (schemaRegistryResponse, error) {
	log.Printf("Processing error response with status code %d", statusCode)

	var errorResponse struct {
		ErrorCode int    `json:"error_code"`
		Message   string `json:"message"`
	}

	if err := json.Unmarshal(body, &errorResponse); err != nil {
		log.Printf("Error parsing error response: %v, body: %s", err, string(body))
		// For parse errors, we return nil for IsCompatible since we couldn't determine compatibility
		return schemaRegistryResponse{
			IsCompatible: nil,
			ErrorCode:    statusCode,
			Message:      fmt.Sprintf("error parsing error response: %v", err),
			HttpStatus:   statusCode,
		}, nil
	}

	log.Printf("Error response parsed - ErrorCode: %d, Message: %s",
		errorResponse.ErrorCode,
		errorResponse.Message)

	// Ensure message has a non-empty string value
	if errorResponse.Message == "" {
		errorResponse.Message = "None"
	}

	// For expected error responses, we set IsCompatible to false
	falseVal := false
	return schemaRegistryResponse{
		IsCompatible: &falseVal,
		ErrorCode:    errorResponse.ErrorCode,
		Message:      errorResponse.Message,
		HttpStatus:   statusCode,
	}, nil
}

// handleSuccessResponse processes successful responses from the schema registry
func handleSuccessResponse(body []byte, statusCode int) (schemaRegistryResponse, error) {
	log.Printf("Processing success response with status code %d", statusCode)

	var result struct {
		IsCompatible bool `json:"is_compatible"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Error parsing success response: %v, body: %s", err, string(body))
		// For parse errors, we return nil for IsCompatible since we couldn't determine compatibility
		return schemaRegistryResponse{
			IsCompatible: nil,
			ErrorCode:    statusCode,
			Message:      fmt.Sprintf("error parsing response: %v", err),
			HttpStatus:   statusCode,
		}, nil
	}

	log.Printf("Success response parsed - IsCompatible: %v", result.IsCompatible)
	// Return a pointer to the boolean value with default message
	return schemaRegistryResponse{
		IsCompatible: &result.IsCompatible,
		ErrorCode:    0,
		Message:      "None", // Default message for success responses
		HttpStatus:   statusCode,
	}, nil
}

// processResponse handles the response based on its status code
func processResponse(body []byte, statusCode int) (schemaRegistryResponse, error) {
	log.Printf("Processing response with status code %d", statusCode)

	switch statusCode {
	case http.StatusNotFound, http.StatusUnprocessableEntity:
		log.Printf("Handling error response (status %d)", statusCode)
		return handleErrorResponse(body, statusCode)
	case http.StatusInternalServerError:
		log.Printf("Handling internal server error (status %d)", statusCode)
		// For internal server errors, set IsCompatible to nil since we couldn't determine compatibility
		return schemaRegistryResponse{
			IsCompatible: nil,
			ErrorCode:    statusCode,
			Message:      "internal server error - please try again later",
			HttpStatus:   statusCode,
		}, nil
	case http.StatusOK:
		log.Printf("Handling success response (status %d)", statusCode)
		return handleSuccessResponse(body, statusCode)
	default:
		log.Printf("Handling unexpected status code %d", statusCode)
		// For unexpected status codes, set IsCompatible to nil since we couldn't determine compatibility
		return schemaRegistryResponse{
			IsCompatible: nil,
			ErrorCode:    statusCode,
			Message:      fmt.Sprintf("unexpected status code: %d, response: %s", statusCode, string(body)),
			HttpStatus:   statusCode,
		}, nil
	}
}

// helper function to transform JSON to Schema Registry format
func transformJSONToSchemaFormat(jsonStr string) (string, error) {
	// First validate that the input is valid JSON
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
	log.Printf("Transformed JSON returned by transformJSONToSchemaFormat: %s", string(result))
	return string(result), nil
}

func setDefault(s *schemaRegistryResponse) {
	s.IsCompatible = nil
	s.ErrorCode = 0
	s.Message = ""
	s.HttpStatus = 0
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, msg string) {
	// Set content type header
	w.Header().Set("Content-Type", "application/json")

	// Set the status code
	w.WriteHeader(statusCode)

	// Create response structure
	response := map[string]interface{}{
		"isValid":    statusCode >= 200 && statusCode < 300,
		"httpStatus": statusCode,
		"message":    msg,
	}
	// Encode and send the response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// If encoding fails, send a simple error
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
