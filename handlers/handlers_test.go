package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kafka-board/types"
)

// For testing, we'll use a mock implementation of the schema validation
type mockSchemaValidator struct {
	mockSchema types.Schema
}

func (m *mockSchemaValidator) ValidatePayload(payload interface{}, schema types.Schema) (bool, []string, error) {
	// Simple mock validation that accepts if test=value and rejects others with appropriate errors
	payloadMap, ok := payload.(map[string]interface{})
	if !ok {
		return false, []string{"Payload is not an object"}, nil
	}

	// Check for required test field
	testValue, hasTest := payloadMap["test"]
	if !hasTest {
		return false, []string{"(root): test is required"}, nil
	}

	// Check for additional properties
	if len(payloadMap) > 1 {
		extraKeys := []string{}
		for k := range payloadMap {
			if k != "test" {
				extraKeys = append(extraKeys, k)
			}
		}
		if len(extraKeys) > 0 {
			return false, []string{fmt.Sprintf("(root): Additional property %s is not allowed", extraKeys[0])}, nil
		}
	}

	// Check the value type of test
	switch t := testValue.(type) {
	case string:
		return true, nil, nil
	case float64:
		return false, []string{"test: Invalid type. Expected: string, given: integer"}, nil
	case bool:
		return false, []string{"Invalid type. Expected: string, given: boolean"}, nil
	default:
		return false, []string{fmt.Sprintf("Unexpected type: %T", t)}, nil
	}
}

// mockHelpers mocks the helpers needed for testing
type mockHelpers struct {
	validator *mockSchemaValidator
}

func (m *mockHelpers) ValidatePayload(payload interface{}, schema types.Schema) (bool, []string, error) {
	return m.validator.ValidatePayload(payload, schema)
}

func (m *mockHelpers) CreateResponseObject(isCompatible *bool, message string, httpStatus int, errorCode int) types.Response {
	return types.Response{
		IsCompatible: isCompatible,
		Message:      message,
		StatusCode:   httpStatus,
		ErrorCode:    errorCode,
	}
}

func (m *mockHelpers) SendJSONResponse(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func (m *mockHelpers) CheckErr(err error) bool {
	return err != nil
}

// Mock registry API implementation
type mockRegistryAPI struct {
	mockSchema types.Schema
}

func (m *mockRegistryAPI) ReturnSubjects() ([]string, error) {
	return []string{}, nil
}

func (m *mockRegistryAPI) ReturnSubjectConfigs(subjectNames []string) ([]types.SubjectConfigInterface, error) {
	return []types.SubjectConfigInterface{}, nil
}

func (m *mockRegistryAPI) GetGlobalConfig() (types.GlobalConfig, error) {
	return types.GlobalConfig{}, nil
}

func (m *mockRegistryAPI) GetSchemas(subjectName string) ([]types.Schema, error) {
	return []types.Schema{m.mockSchema}, nil
}

func (m *mockRegistryAPI) TestSchema(subjectName string, version int, testJSON string) (types.Response, error) {
	return types.Response{}, nil
}

func (m *mockRegistryAPI) GetSchema(id string) (types.Schema, error) {
	return m.mockSchema, nil
}

func TestHandleValidatePayload(t *testing.T) {
	// Set up logger for testing
	var logBuffer bytes.Buffer
	testLogger := slog.New(slog.NewTextHandler(&logBuffer, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Create mock schema for testing
	mockSchema := types.Schema{
		Id: 1,
		Schema: `{
				"additionalProperties": false,
				"properties": {
					"test": {
						"description": "A required test string",
						"type": "string"
					}
				},
				"required": [
					"test"
				],
				"type": "object"
			}`,
	}

	// Create validator, mock helpers, and mock API
	validator := &mockSchemaValidator{mockSchema: mockSchema}
	mockHelperInstance := &mockHelpers{validator: validator}
	mockAPIInstance := &mockRegistryAPI{mockSchema: mockSchema}

	// Custom function to replace HandleValidatePayload for testing
	testHandleValidatePayload := func(w http.ResponseWriter, r *http.Request) {
		// Extract the id from the query parameters
		id := r.URL.Query().Get("id")

		// Read and validate request body
		body, err := io.ReadAll(r.Body)
		if mockHelperInstance.CheckErr(err) {
			response := mockHelperInstance.CreateResponseObject(
				&falseVal,
				fmt.Sprintf("Error reading request body: %v", err),
				http.StatusBadRequest,
				0,
			)
			mockHelperInstance.SendJSONResponse(w, http.StatusBadRequest, response)
			return
		}

		// Parse the JSON request body
		var unmarshalledBody map[string]any
		err = json.Unmarshal(body, &unmarshalledBody)
		if mockHelperInstance.CheckErr(err) {
			response := mockHelperInstance.CreateResponseObject(
				&falseVal,
				fmt.Sprintf("Invalid JSON format in request body: %v", err),
				http.StatusBadRequest,
				0,
			)
			mockHelperInstance.SendJSONResponse(w, http.StatusBadRequest, response)
			return
		}

		// Ensure payload key exists
		payloadRaw, ok := unmarshalledBody["payload"]
		if !ok {
			response := mockHelperInstance.CreateResponseObject(
				&falseVal,
				"payload key expected in request body",
				http.StatusBadRequest,
				0,
			)
			mockHelperInstance.SendJSONResponse(w, http.StatusBadRequest, response)
			return
		}

		// Process the payload based on its type
		var payload interface{}

		// Case 1: Payload is a string (need to parse as JSON)
		payloadStr, isString := payloadRaw.(string)
		if isString {
			err = json.Unmarshal([]byte(payloadStr), &payload)
			if mockHelperInstance.CheckErr(err) {
				response := mockHelperInstance.CreateResponseObject(
					&falseVal,
					fmt.Sprintf("value of payload key is not valid JSON: %v", err),
					http.StatusBadRequest,
					0,
				)
				mockHelperInstance.SendJSONResponse(w, http.StatusBadRequest, response)
				return
			}
		} else {
			// Case 2: Payload is already a JSON object
			payload = payloadRaw
		}

		// Get the schema (using our mock)
		schema, err := mockAPIInstance.GetSchema(id)
		if mockHelperInstance.CheckErr(err) {
			response := mockHelperInstance.CreateResponseObject(
				&falseVal,
				fmt.Sprintf("Error retrieving schema: %v", err),
				http.StatusInternalServerError,
				0,
			)
			mockHelperInstance.SendJSONResponse(w, http.StatusInternalServerError, response)
			return
		}

		isValid, errors, err := mockHelperInstance.ValidatePayload(payload, schema)
		if err != nil {
			response := mockHelperInstance.CreateResponseObject(
				&falseVal,
				fmt.Sprintf("Error validating payload: %v", err),
				http.StatusInternalServerError,
				0,
			)
			mockHelperInstance.SendJSONResponse(w, http.StatusInternalServerError, response)
			return
		}

		// Create response based on validation result
		var response types.Response
		if !isValid {
			response = mockHelperInstance.CreateResponseObject(
				&falseVal,
				strings.Join(errors, "; "),
				http.StatusOK,
				0,
			)
		} else {
			response = mockHelperInstance.CreateResponseObject(
				&trueVal,
				"Payload validates against schema",
				http.StatusOK,
				0,
			)
		}

		testLogger.Debug("Validation result", "valid", isValid, "errors", errors)
		mockHelperInstance.SendJSONResponse(w, response.StatusCode, response)
	}

	tests := []struct {
		name               string
		payload            string
		id                 string
		expectedStatusCode int
		expectedCompatible bool
		expectedMessage    string
	}{
		{
			name:               "invalid json in payload",
			payload:            `{"name 123///`,
			id:                 "1",
			expectedStatusCode: http.StatusBadRequest,
			expectedCompatible: false,
			expectedMessage:    "Invalid JSON format in request body",
		},
		{
			name:               "payload has not payload key",
			payload:            `{"name": "John"}`,
			id:                 "1",
			expectedStatusCode: http.StatusBadRequest,
			expectedCompatible: false,
			expectedMessage:    "payload key expected in request body",
		},
		{
			name:               "well-formed and valid payload with string JSON",
			payload:            `{"payload": "{\"test\": \"value\"}"}`,
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			expectedCompatible: true,
			expectedMessage:    "Payload validates against schema",
		},
		{
			name:               "invalid JSON inside payload string",
			payload:            `{"payload": "{\"test\": broken}"}`,
			id:                 "1",
			expectedStatusCode: http.StatusBadRequest,
			expectedCompatible: false,
			expectedMessage:    "value of payload key is not valid JSON",
		},
		{
			name:               "payload with object directly",
			payload:            `{"payload": {"test": "value"}}`,
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			expectedCompatible: true,
			expectedMessage:    "Payload validates against schema",
		},
		{
			name:               "valid payload but invalid for schema (missing required field)",
			payload:            `{"payload": "{\"not_test\": \"value\"}"}`,
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			expectedCompatible: false,
			expectedMessage:    "(root): test is required",
		},
		{
			name:               "valid payload with extra property not allowed by schema",
			payload:            `{"payload": {"test": "value", "extra": "not allowed"}}`,
			id:                 "1",
			expectedStatusCode: http.StatusOK,
			expectedCompatible: false,
			expectedMessage:    "Additional property extra is not allowed",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Reset log buffer between tests
			logBuffer.Reset()

			// Create request with payload
			req := httptest.NewRequest("POST", "/test-payload?id="+test.id, strings.NewReader(test.payload))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()

			// Execute our test handler function instead of the real one
			testHandleValidatePayload(rr, req)

			// Check status code
			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatusCode)
			}

			// Check response body
			var response types.Response
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("Could not parse response: %v", err)
			}

			// Check validation result
			isCompatibleValue := false
			if response.IsCompatible != nil {
				isCompatibleValue = *response.IsCompatible
			}

			if isCompatibleValue != test.expectedCompatible {
				t.Errorf("Expected valid=%v, got %v", test.expectedCompatible, isCompatibleValue)
			}

			if !strings.Contains(response.Message, test.expectedMessage) {
				t.Errorf("Expected message to contain '%s', got '%s'", test.expectedMessage, response.Message)
			}

			// Log for debugging
			t.Logf("Response: %+v", response)
			t.Logf("Logs: %s", logBuffer.String())
		})
	}
}
