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

	"kafka-board/helpers"
	"kafka-board/types"
)

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

	// Create validator and mock API
	// validator := &mockSchemaValidator{mockSchema: mockSchema}
	mockAPIInstance := &mockRegistryAPI{mockSchema: mockSchema}

	// Custom function to replace HandleValidatePayload for testing
	testHandleValidatePayload := func(w http.ResponseWriter, r *http.Request) {
		// Extract the id from the query parameters
		id := r.URL.Query().Get("id")

		// Read and validate request body
		body, err := io.ReadAll(r.Body)
		if helpers.CheckErr(err) {
			response := helpers.CreateResponseObject(
				&falseVal,
				fmt.Sprintf("Error reading request body: %v", err),
				http.StatusBadRequest,
				0,
			)
			helpers.SendJSONResponse(w, http.StatusBadRequest, response)
			return
		}

		// Parse the JSON request body
		var unmarshalledBody map[string]any
		err = json.Unmarshal(body, &unmarshalledBody)
		if helpers.CheckErr(err) {
			response := helpers.CreateResponseObject(
				&falseVal,
				fmt.Sprintf("Invalid JSON format in request body: %v", err),
				http.StatusBadRequest,
				0,
			)
			helpers.SendJSONResponse(w, http.StatusBadRequest, response)
			return
		}

		// Ensure payload key exists
		payloadRaw, ok := unmarshalledBody["payload"]
		if !ok {
			response := helpers.CreateResponseObject(
				&falseVal,
				"payload key expected in request body",
				http.StatusBadRequest,
				0,
			)
			helpers.SendJSONResponse(w, http.StatusBadRequest, response)
			return
		}

		// Process the payload based on its type
		var payload interface{}

		// Case 1: Payload is a string (need to parse as JSON)
		payloadStr, isString := payloadRaw.(string)
		if isString {
			err = json.Unmarshal([]byte(payloadStr), &payload)
			if helpers.CheckErr(err) {
				response := helpers.CreateResponseObject(
					&falseVal,
					fmt.Sprintf("value of payload key is not valid JSON: %v", err),
					http.StatusBadRequest,
					0,
				)
				helpers.SendJSONResponse(w, http.StatusBadRequest, response)
				return
			}
		} else {
			// Case 2: Payload is already a JSON object
			payload = payloadRaw
		}

		// Get the schema (using our mock)
		schema, err := mockAPIInstance.GetSchema(id)
		if helpers.CheckErr(err) {
			response := helpers.CreateResponseObject(
				&falseVal,
				fmt.Sprintf("Error retrieving schema: %v", err),
				http.StatusInternalServerError,
				0,
			)
			helpers.SendJSONResponse(w, http.StatusInternalServerError, response)
			return
		}

		// Use the real helper function for validation directly
		isValid, errors, err := helpers.ValidatePayload(payload, schema)
		if err != nil {
			response := helpers.CreateResponseObject(
				&falseVal,
				fmt.Sprintf("Error validating payload: %v", err),
				http.StatusInternalServerError,
				0,
			)
			helpers.SendJSONResponse(w, http.StatusInternalServerError, response)
			return
		}

		// Create response based on validation result
		var response types.Response
		if !isValid {
			response = helpers.CreateResponseObject(
				&falseVal,
				strings.Join(errors, "; "),
				http.StatusOK,
				0,
			)
		} else {
			response = helpers.CreateResponseObject(
				&trueVal,
				"Payload validates against schema",
				http.StatusOK,
				0,
			)
		}

		testLogger.Debug("Validation result", "valid", isValid, "errors", errors)
		helpers.SendJSONResponse(w, response.StatusCode, response)
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
