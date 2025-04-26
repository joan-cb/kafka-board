package handlers

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kafka-board/types"
)

func TestHandleValidatePayload(t *testing.T) {
	// Set up logger for testing
	t.Parallel()
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

	testHandler := ReturnHandler(testLogger, &mockRegistryAPI{mockSchema: mockSchema})

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
			testHandler.HandleValidatePayload(rr, req)

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
