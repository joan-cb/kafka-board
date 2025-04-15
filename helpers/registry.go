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

// isEmptyJSON checks if a parsed JSON value is empty (empty object, array, string, or null)
func (helper *Helpers) isEmptyJSON(jsonObj interface{}) bool {
	switch v := jsonObj.(type) {
	case map[string]interface{}:
		return len(v) == 0 // empty object: {}
	case []interface{}:
		return len(v) == 0 // empty array: []
	case string:
		return v == "" // empty string
	case nil:
		return true // null
	default:
		return false // other types (numbers, booleans) are not considered empty
	}
}

// TransformJSONToSchemaFormat takes a JSON string and wraps it in the Schema Registry format.
// It validates that the input is valid JSON before creating the wrapper structure.
// This function is used to prepare JSON schemas for schema registry compatibility testing.
//
// Parameters:
//   - jsonStr: A string containing valid JSON to be wrapped
//
// Returns:
//   - string: The JSON string in Schema Registry format
//   - error: An error if the JSON is invalid or if marshaling fails
func (helper *Helpers) TransformJSONToSchemaFormat(jsonStr string) (string, error) {
	// First validate the JSON by attempting to unmarshal it
	var jsonObj interface{}

	if isEmpty := helper.isEmptyJSON(jsonStr); isEmpty {
		helper.logger.Debug("TransformJSONToSchemaFormat - empty JSON provided",
			"jsonStr", jsonStr)
		return "", fmt.Errorf("empty JSON is not allowed")
	}

	helper.logger.Debug("TransformJSONToSchemaFormat - validating JSON input",
		"jsonStr", jsonStr)
	if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
		return "", fmt.Errorf("invalid JSON input: %v", err)
	}

	// Check for empty JSON objects or arrays
	if helper.isEmptyJSON(jsonObj) {
		helper.logger.Debug("TransformJSONToSchemaFormat - empty JSON provided")
		return "", fmt.Errorf("empty JSON is not allowed")
	}

	// Create the schema registry format wrapper
	schemaRegistryFormat := SchemaFormat{
		Schema:     jsonStr,
		SchemaType: "JSON",
	}

	// Marshal to JSON then unmarshal to map to check keys with key, ok idiom
	tmpJSON, err := json.Marshal(schemaRegistryFormat)
	if err != nil {
		helper.logger.Debug("TransformJSONToSchemaFormat - error formatting schema",
			"error", err)
		return "", fmt.Errorf("error formatting schema: %v", err)
	}

	var schemaMap map[string]interface{}
	if err := json.Unmarshal(tmpJSON, &schemaMap); err != nil {
		helper.logger.Debug("TransformJSONToSchemaFormat - error validating schema format",
			"error", err)
		return "", fmt.Errorf("error validating schema format: %v", err)
	}

	// Use key, ok idiom to check for required fields
	if _, ok := schemaMap["schema"]; !ok {
		helper.logger.Debug("TransformJSONToSchemaFormat - missing required field: schema",
			"error", err)
		return "", fmt.Errorf("missing required field: schema")
	}

	if _, ok := schemaMap["schemaType"]; !ok {
		helper.logger.Debug("TransformJSONToSchemaFormat - missing required field: schemaType",
			"error", err)
		return "", fmt.Errorf("missing required field: schemaType")
	}

	// Return the marshalled result
	return string(tmpJSON), nil
}

// ProcessResponse processes HTTP responses from the Schema Registry.
// It handles both success and error cases, extracting the relevant information from the response body.
//
// Parameters:
//   - body: The response body as a byte array
//   - statusCode: The HTTP status code from the response
//
// Returns:
//   - Response: A structured response with compatibility information or error details
//   - error: An error if processing fails or the registry reported an error
func (helper *Helpers) ProcessResponse(body []byte, statusCode int) (types.Response, error) {
	// Handle different status codes appropriately
	switch {
	case statusCode == http.StatusOK:
		// Success case - try to parse compatibility result
		var result struct {
			IsCompatible bool `json:"is_compatible"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			errMsg := fmt.Sprintf("error parsing compatibility response: %v", err)
			helper.logger.Debug("ProcessResponse - error parsing compatibility response",
				"error", errMsg)
			return CreateResponseObject(
				nil,
				errMsg,
				statusCode,
				statusCode,
			), fmt.Errorf("%s", errMsg)
		}

		helper.logger.Debug("ProcessResponse - success with valid compatibility result",
			"isCompatible", result.IsCompatible)

		// Success with valid compatibility result
		return CreateResponseObject(
			&result.IsCompatible,
			"None", // Default message for success
			statusCode,
			0,
		), nil

	case statusCode == http.StatusNotFound || statusCode == http.StatusUnprocessableEntity:
		// Expected error responses from registry - try to parse error details
		var errorResponse struct {
			ErrorCode int    `json:"error_code"`
			Message   string `json:"message"`
		}

		if err := json.Unmarshal(body, &errorResponse); err != nil {
			errMsg := fmt.Sprintf("error parsing error response: %v", err)

			helper.logger.Debug("ProcessResponse - error parsing error response",
				"error", errMsg)

			return CreateResponseObject(
				nil,
				errMsg,
				statusCode,
				statusCode,
			), fmt.Errorf("%s", errMsg)
		}

		// Ensure message has a value
		if errorResponse.Message == "" {
			errorResponse.Message = "Registry returned error with no message"
		}

		// Create descriptive error with registry details
		err := fmt.Errorf("registry error: %s (code: %d, status: %d)",
			errorResponse.Message,
			errorResponse.ErrorCode,
			statusCode)

		// Set IsCompatible to false for expected errors
		falseVal := false

		helper.logger.Debug("ProcessResponse - error with expected error response",
			"error", err)

		return CreateResponseObject(
			&falseVal,
			errorResponse.Message,
			statusCode,
			errorResponse.ErrorCode,
		), err

	case statusCode == http.StatusInternalServerError:
		// Internal server error in registry
		helper.logger.Debug("ProcessResponse - error with internal server error")

		return CreateResponseObject(
			nil,
			"internal server error - attempt to validate schema failed",
			statusCode,
			0,
		), fmt.Errorf("registry internal server error: status code %d", statusCode)

	default:
		// Unexpected status codes
		errMsg := fmt.Sprintf("unexpected status code: %d, response: %s", statusCode, string(body))

		helper.logger.Debug("ProcessResponse - error with unexpected status code",
			"error", errMsg)

		return CreateResponseObject(
			nil,
			errMsg,
			statusCode,
			statusCode,
		), fmt.Errorf("registry returned unexpected status: %s", errMsg)
	}
}
