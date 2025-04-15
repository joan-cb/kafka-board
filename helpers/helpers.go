package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"kafka-board/types"
)

// TO DO: Add STRUCTURED logger to helpers

var baseRegistryURL = "http://schema-registry:8081"

// CreateTestSchemaRequest creates a new HTTP request for testing schema compatibility
func CreateTestSchemaRequest(subjectName string, version int, testJSON string) (*http.Request, error) {
	requestURL := baseRegistryURL + "/compatibility/subjects/" + subjectName + "/versions/" + strconv.Itoa(version)
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(testJSON)))
	if err != nil {

		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// MakeHTTPRequest executes an HTTP request and returns the response
func MakeHTTPRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		return nil, fmt.Errorf("error making request: %v", err)
	}

	return resp, nil
}

// ReadResponseBody reads and returns the response body
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {

		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return body, nil
}

// ProcessCompatibilityResponse handles the response based on its status code
func ProcessCompatibilityResponse(body []byte, statusCode int) (types.SchemaRegistryResponse, error) {
	switch statusCode {
	case http.StatusNotFound, http.StatusUnprocessableEntity:
		return handleErrorResponse(body, statusCode)
	case http.StatusInternalServerError:
		// For internal server errors, set IsCompatible to nil since we couldn't determine compatibility
		return types.SchemaRegistryResponse{
			IsCompatible: nil,
			ErrorCode:    0,
			Message:      "internal server error - please try again later",
			StatusCode:   statusCode,
		}, nil
	case http.StatusOK:
		return handleSuccessResponse(body, statusCode)
	default:
		// For unexpected status codes, set IsCompatible to nil since we couldn't determine compatibility
		return types.SchemaRegistryResponse{
			IsCompatible: nil,
			ErrorCode:    statusCode,
			Message:      fmt.Sprintf("unexpected status code: %d, response: %s", statusCode, string(body)),
			StatusCode:   statusCode,
		}, nil
	}
}

// handleErrorResponse processes error responses from the schema registry
func handleErrorResponse(body []byte, statusCode int) (types.SchemaRegistryResponse, error) {
	var errorResponse struct {
		ErrorCode int    `json:"error_code"`
		Message   string `json:"message"`
	}
	if err := json.Unmarshal(body, &errorResponse); err != nil {

		resp := CreateSchemaRegistryResponse(nil, fmt.Sprintf("error parsing error response: %v", err), statusCode, statusCode)
		return resp, nil
	}

	// Ensure message has a non-empty string value
	if errorResponse.Message == "" {
		errorResponse.Message = "None"
	}
	// For expected error responses, we set IsCompatible to false
	falseVal := false
	resp := CreateSchemaRegistryResponse(
		&falseVal,
		errorResponse.Message,
		statusCode,
		errorResponse.ErrorCode,
	)

	return resp, nil
}

// handleSuccessResponse processes successful responses from the schema registry
func handleSuccessResponse(body []byte, statusCode int) (types.SchemaRegistryResponse, error) {
	var result struct {
		IsCompatible bool `json:"is_compatible"`
	}

	if err := json.Unmarshal(body, &result); err != nil {

		resp := CreateSchemaRegistryResponse(nil, fmt.Sprintf("error parsing response: %v", err), statusCode, statusCode)
		return resp, nil
	}

	// Return a pointer to the boolean value with default message
	resp := CreateSchemaRegistryResponse(&result.IsCompatible, "None", statusCode, 0)

	return resp, nil
}

// TransformJSONToSchemaFormat takes a JSON string and wraps it in the Schema Registry format
func TransformJSONToSchemaFormat(jsonStr string) (string, error) {
	// Unmarshal the JSON string into an interface{} and not a [string]interface{} to allow for more flexible JSON input e.g. arrays or objects or primitives
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
func CreatePayloadResponse(isValid bool, message string, statusCode int) types.PayloadTestResponse {
	return types.PayloadTestResponse{
		IsCompatible: isValid,
		Message:      message,
		StatusCode:   statusCode,
	}
}

// CreateSchemaRegistryResponse creates a schemaRegistryResponse struct instance
func CreateSchemaRegistryResponse(isCompatible *bool, message string, httpStatus int, errorCode int) types.SchemaRegistryResponse {
	return types.SchemaRegistryResponse{
		IsCompatible: isCompatible,
		Message:      message,
		StatusCode:   httpStatus,
		ErrorCode:    errorCode,
	}
}

// CheckErr is a helper function to check if an error is present
func CheckErr(e error) bool {
	return e != nil
}

func SetupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: getLogLevel(),
	}))
}

func GetLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func GetServerAddress() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9080"
	}
	return ":" + port
}
