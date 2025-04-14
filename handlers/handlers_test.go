package handlers

// import (
// 	"encoding/json"
// 	"log/slog"
// 	"net/http"
// 	"net/http/httptest"
// 	"os"
// 	"strings"
// 	"testing"
// )

// func TestHandleValidatePayload(t *testing.T) {
// 	// Set up logger to discard logs to avoid panics
// 	logggerHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
// 		Level: slog.LevelDebug,
// 	})
// 	logger = slog.New(logggerHandler)
// 	logger = slog.New(logggerHandler)

// 	mockSchema := Schema{
// 		Id: 1,
// 		Schema: `{
// 				"additionalProperties": false,
// 				"properties": {
// 					"test": {
// 						"description": "A required test string",
// 						"type": "string"
// 					}
// 				},
// 				"required": [
// 					"test"
// 				],
// 				"type": "object"
// 			}`,
// 	}

// 	testHandler := returnHandler(&mockHandler{
// 		mockSchema: mockSchema,
// 	})

// 	tests := []struct {
// 		name               string
// 		payload            string
// 		id                 string
// 		expectedStatusCode int
// 		expectedCompatible bool
// 		expectedMessage    string
// 	}{
// 		{
// 			name:               "invalid json in payload",
// 			payload:            `{"name 123///`,
// 			id:                 "1",
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedCompatible: false,
// 			expectedMessage:    "Invalid JSON format in request body",
// 		},
// 		{
// 			name:               "payload has not payload key",
// 			payload:            `{"name": "John"}`,
// 			id:                 "1",
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedCompatible: false,
// 			expectedMessage:    "payload key expected in request body",
// 		},
// 		{
// 			name:               "well-formed and valid payload with string JSON",
// 			payload:            `{"payload": "{\"test\": \"value\"}"}`,
// 			id:                 "1",
// 			expectedStatusCode: http.StatusOK,
// 			expectedCompatible: true,
// 			expectedMessage:    "Payload validates against schema",
// 		},
// 		{
// 			name:               "invalid JSON inside payload string",
// 			payload:            `{"payload": "{\"test\": broken}"}`,
// 			id:                 "1",
// 			expectedStatusCode: http.StatusBadRequest,
// 			expectedCompatible: false,
// 			expectedMessage:    "value of payload key is not valid JSON: invalid character 'b' looking for beginning of value",
// 		},
// 		{
// 			name:               "payload with object directly",
// 			payload:            `{"payload": {"test": "value"}}`,
// 			id:                 "1",
// 			expectedStatusCode: http.StatusOK,
// 			expectedCompatible: true,
// 			expectedMessage:    "Payload validates against schema",
// 		},
// 		{
// 			name:               "valid payload with string JSON but invalid for schema",
// 			payload:            `{"payload": "{\"test_invalid\": \"value\"}"}`,
// 			id:                 "1",
// 			expectedStatusCode: http.StatusOK,
// 			expectedCompatible: false,
// 			expectedMessage:    "(root): test is required",
// 		},
// 		{
// 			name:               "valid payload with string JSON but invalid for schema - wrong type",
// 			payload:            `{"payload": "{\"test\": 1}"}`,
// 			id:                 "1",
// 			expectedStatusCode: http.StatusOK,
// 			expectedCompatible: false,
// 			expectedMessage:    "test: Invalid type. Expected: string, given: integer",
// 		},
// 		{
// 			name:               "valid payload with string JSON but invalid for schema - wrong type",
// 			payload:            `{"payload": "{\"test\": false}"}`,
// 			id:                 "1",
// 			expectedStatusCode: http.StatusOK,
// 			expectedCompatible: false,
// 			expectedMessage:    "Invalid type. Expected: string, given: boolean",
// 		},
// 		{
// 			name:               "valid payload with string JSON but has an additional property - which is not allowed by the schema",
// 			payload:            `{"payload": {"test": "str", "extra": "value"}}`,
// 			id:                 "1",
// 			expectedStatusCode: http.StatusOK,
// 			expectedCompatible: false, // set to  true to validate the message
// 			expectedMessage:    "(root): Additional property extra is not allowed",
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			// Create request with payload
// 			t.Parallel()

// 			// Use the string payload directly
// 			req := httptest.NewRequest("POST", "/test-payload?id="+test.id, strings.NewReader(test.payload))
// 			req.Header.Set("Content-Type", "application/json")

// 			// Create response recorder
// 			rr := httptest.NewRecorder()

// 			// Execute the handler
// 			testHandler.handleValidatePayload(rr, req)

// 			// Check status code
// 			if status := rr.Code; status != test.expectedStatusCode {
// 				t.Errorf("handler returned wrong status code: got %v want %v", status, test.expectedStatusCode)
// 			}

// 			// Check response body
// 			var response payloadTestResponse
// 			err := json.Unmarshal(rr.Body.Bytes(), &response)
// 			if err != nil {
// 				t.Fatalf("Could not parse response: %v", err)
// 			}

// 			// Check validation result
// 			if response.IsCompatible != test.expectedCompatible {
// 				t.Errorf("Expected valid=%v, got %v", test.expectedCompatible, response.IsCompatible)
// 			}
// 			if !strings.Contains(response.Message, test.expectedMessage) {
// 				t.Errorf("Expected message to contain '%s', got '%s'", test.expectedMessage, response.Message)
// 			}

// 			// Log for debugging
// 			t.Logf("Response: %+v", response)
// 		})
// 	}
// }
