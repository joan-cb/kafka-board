package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"kafka-board/types"
)

// SchemaFormat represents the structure required by the Schema Registry API
type SchemaFormat struct {
	Schema     string `json:"schema"`
	SchemaType string `json:"schemaType"`
}

// TransformJSONToSchemaFormat takes a JSON string and wraps it in the Schema Registry format.
// It validates that the input is valid JSON before creating the wrapper structure.
// This function is used to prepare JSON schemas for compatibility testing.
//
// Parameters:
//   - jsonStr: A string containing valid JSON to be wrapped
//
// Returns:
//   - string: The JSON string in Schema Registry format
//   - error: An error if the JSON is invalid or if marshaling fails
func TransformJSONToSchemaFormat(jsonStr string) (string, error) {
	// First validate the JSON by attempting to unmarshal it
	var jsonObj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
		return "", fmt.Errorf("invalid JSON input: %v", err)
	}

	// Create the schema registry format wrapper
	schemaRegistryFormat := SchemaFormat{
		Schema:     jsonStr,
		SchemaType: "JSON",
	}

	// Marshal the wrapper back to JSON
	result, err := json.Marshal(schemaRegistryFormat)
	if err != nil {
		return "", fmt.Errorf("error formatting schema: %v", err)
	}

	return string(result), nil
}

// ProcessCompatibilityResponse handles the Schema Registry API response based on its status code.
// It determines if a schema is compatible and extracts any error information.
//
// Parameters:
//   - body: The response body as a byte array
//   - statusCode: The HTTP status code from the response
//
// Returns:
//   - SchemaRegistryResponse: A structured response with compatibility information
//   - error: An error if processing fails
func ProcessCompatibilityResponse(body []byte, statusCode int) (types.Response, error) {
	switch statusCode {
	case http.StatusNotFound, http.StatusUnprocessableEntity:
		// Expected error cases (subject not found, schema parsing failed, etc.)
		return handleErrorResponse(body, statusCode)

	case http.StatusInternalServerError:
		// For internal server errors, we can't determine compatibility
		return types.Response{
			IsCompatible: nil,
			ErrorCode:    0,
			Message:      "internal server error - attempt to validate schema failed",
			StatusCode:   statusCode,
		}, fmt.Errorf("registry internal server error: status code %d", statusCode)

	case http.StatusOK:
		// Success case - check compatibility result
		return handleSuccessResponse(body, statusCode)

	default:
		// Unexpected status codes
		errMsg := fmt.Sprintf("unexpected status code: %d, response: %s", statusCode, string(body))
		return types.Response{
			IsCompatible: nil,
			ErrorCode:    statusCode,
			Message:      errMsg,
			StatusCode:   statusCode,
		}, fmt.Errorf("registry returned unexpected status: %s", errMsg)
	}
}

// handleErrorResponse processes error responses from the schema registry.
// It extracts the error code and message from the response body.
//
// Parameters:
//   - body: The response body as a byte array
//   - statusCode: The HTTP status code from the response
//
// Returns:
//   - SchemaRegistryResponse: A structured response with error information
//   - error: Error information extracted from the response
func handleErrorResponse(body []byte, statusCode int) (types.Response, error) {
	var errorResponse struct {
		ErrorCode int    `json:"error_code"`
		Message   string `json:"message"`
	}

	if err := json.Unmarshal(body, &errorResponse); err != nil {
		// Handle JSON parsing errors in the error response
		errMsg := fmt.Sprintf("error parsing error response: %v", err)
		return CreateResponse(
			nil,
			errMsg,
			statusCode,
			statusCode,
		), fmt.Errorf("%s", errMsg)
	}

	// Ensure message has a non-empty string value
	if errorResponse.Message == "" {
		errorResponse.Message = "Registry returned error with no message"
	}

	// Create a descriptive error that includes the registry's error details
	err := fmt.Errorf("registry error: %s (code: %d, status: %d)",
		errorResponse.Message,
		errorResponse.ErrorCode,
		statusCode)

	// For expected error responses, we set IsCompatible to false
	falseVal := false
	return CreateResponse(
		&falseVal,
		errorResponse.Message,
		statusCode,
		errorResponse.ErrorCode,
	), err
}

// handleSuccessResponse processes successful responses from the schema registry.
// It extracts the compatibility result from the response body.
//
// Parameters:
//   - body: The response body as a byte array
//   - statusCode: The HTTP status code from the response
//
// Returns:
//   - SchemaRegistryResponse: A structured response with compatibility information
//   - error: Error if JSON parsing fails, nil otherwise
func handleSuccessResponse(body []byte, statusCode int) (types.Response, error) {
	var result struct {
		IsCompatible bool `json:"is_compatible"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		// Handle JSON parsing errors in the success response
		errMsg := fmt.Sprintf("error parsing compatibility response: %v", err)
		return CreateResponse(
			nil,
			errMsg,
			statusCode,
			statusCode,
		), fmt.Errorf("%s", errMsg)
	}

	// Return a pointer to the boolean value with default message
	// On success, return nil error
	return CreateResponse(
		&result.IsCompatible,
		"None",
		statusCode,
		0,
	), nil
}
