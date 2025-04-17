package helpers

// import (
// 	"io"
// 	"testing"
// )

// func TestCreateTestSchemaRequest(t *testing.T) {
// 	// Save original baseRegistryURL and restore after test
// 	originalURL := baseRegistryURL
// 	defer func() { baseRegistryURL = originalURL }()

// 	// Set test URL
// 	baseRegistryURL = "http://test-registry:8081"

// 	tests := []struct {
// 		name          string
// 		subject       string
// 		version       int
// 		json          string
// 		wantURL       string
// 		wantErr       bool
// 		msg           string
// 		expectedError string
// 	}{
// 		{
// 			name:    "valid request",
// 			subject: "test-subject",
// 			version: 1,
// 			json:    `{"type":"string"}`,
// 			wantURL: "http://test-registry:8081/compatibility/subjects/test-subject/versions/1",
// 			wantErr: false,
// 			msg:     "valid request",
// 		},
// 		{
// 			name:    "subject with special chars",
// 			subject: "test.subject-name",
// 			version: 2,
// 			json:    `{"type":"object"}`,
// 			wantURL: "http://test-registry:8081/compatibility/subjects/test.subject-name/versions/2",
// 			wantErr: false,
// 			msg:     "subject with special chars",
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			req, err := CreateTestSchemaRequest(test.subject, test.version, test.json)

// 			// Check error result
// 			if (err != nil) != test.wantErr {
// 				t.Errorf("createTestSchemaRequest() error = %v, wantErr %v", err, test.wantErr)
// 				return
// 			}
// 			//returns if err!nil to avoid a cause a nil pointer panic.
// 			if err != nil {

// 				return
// 			}

// 			// Check URL
// 			if req.URL.String() != test.wantURL {
// 				t.Errorf("URL = %v, want %v", req.URL.String(), test.wantURL)
// 			}

// 			// Check headers
// 			if req.Header.Get("Accept") != "application/vnd.schemaregistry.v1+json" {
// 				t.Errorf("Accept header = %v, want %v",
// 					req.Header.Get("Accept"), "application/vnd.schemaregistry.v1+json")
// 			}

// 			if req.Header.Get("Content-Type") != "application/json" {
// 				t.Errorf("Content-Type header = %v, want %v",
// 					req.Header.Get("Content-Type"), "application/json")
// 			}

// 			// Verify body content
// 			bodyBytes, _ := io.ReadAll(req.Body)
// 			if string(bodyBytes) != test.json {
// 				t.Errorf("Body = %v, want %v", string(bodyBytes), test.json)
// 			}
// 		})
// 	}
// }
