package helpers

// func TestIsEmptyJSON(t *testing.T) {
// 	// Create a helpers instance for testing
// 	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
// 	helper := &Helpers{logger: logger}

// 	tests := []struct {
// 		name     string
// 		input    interface{}
// 		expected bool
// 	}{
// 		{
// 			name:     "empty object",
// 			input:    map[string]interface{}{},
// 			expected: true,
// 		},
// 		{
// 			name:     "non-empty object",
// 			input:    map[string]interface{}{"key": "value"},
// 			expected: false,
// 		},
// 		{
// 			name:     "empty array",
// 			input:    []interface{}{},
// 			expected: true,
// 		},
// 		{
// 			name:     "non-empty array",
// 			input:    []interface{}{1, 2, 3},
// 			expected: false,
// 		},
// 		{
// 			name:     "empty string",
// 			input:    "",
// 			expected: true,
// 		},
// 		{
// 			name:     "non-empty string",
// 			input:    "hello",
// 			expected: false,
// 		},
// 		{
// 			name:     "null",
// 			input:    nil,
// 			expected: true,
// 		},
// 		{
// 			name:     "integer",
// 			input:    42,
// 			expected: false,
// 		},
// 		{
// 			name:     "boolean true",
// 			input:    true,
// 			expected: false,
// 		},
// 		{
// 			name:     "boolean false",
// 			input:    false,
// 			expected: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result := helper.isEmptyJSON(tt.input)
// 			if result != tt.expected {
// 				t.Errorf("isEmptyJSON() = %v, want %v", result, tt.expected)
// 			}
// 		})
// 	}
// }

// func TestTransformJSONToSchemaFormat(t *testing.T) {
// 	// Create a helpers instance for testing
// 	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
// 	helper := &Helpers{logger: logger}

// 	tests := []struct {
// 		name        string
// 		input       string
// 		wantErr     bool
// 		expectedErr string
// 	}{
// 		{
// 			name:    "valid simple JSON object",
// 			input:   `{"field1": "value1", "field2": 42}`,
// 			wantErr: false,
// 		},
// 		{
// 			name:    "valid complex JSON object",
// 			input:   `{"type": "object", "properties": {"name": {"type": "string"}, "age": {"type": "integer"}}}`,
// 			wantErr: false,
// 		},
// 		{
// 			name:    "valid JSON array",
// 			input:   `[1, 2, 3, 4]`,
// 			wantErr: false,
// 		},
// 		{
// 			name:        "invalid JSON - missing closing brace",
// 			input:       `{"field1": "value1"`,
// 			wantErr:     true,
// 			expectedErr: "invalid JSON input",
// 		},
// 		{
// 			name:        "invalid JSON - trailing comma",
// 			input:       `{"field1": "value1", }`,
// 			wantErr:     true,
// 			expectedErr: "invalid JSON input",
// 		},
// 		{
// 			name:        "empty string",
// 			input:       "",
// 			wantErr:     true,
// 			expectedErr: "empty JSON is not allowed",
// 		},
// 		{
// 			name:        "empty JSON object",
// 			input:       "{}",
// 			wantErr:     true,
// 			expectedErr: "empty JSON is not allowed",
// 		},
// 		{
// 			name:        "empty JSON array",
// 			input:       "[]",
// 			wantErr:     true,
// 			expectedErr: "empty JSON is not allowed",
// 		},
// 		{
// 			name:        "null JSON",
// 			input:       "null",
// 			wantErr:     true,
// 			expectedErr: "empty JSON is not allowed",
// 		},
// 		{
// 			name:        "empty string value",
// 			input:       `""`,
// 			wantErr:     true,
// 			expectedErr: "empty JSON is not allowed",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := helper.TransformJSONToSchemaFormat(tt.input)

// 			// Check error expectations
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("TransformJSONToSchemaFormat() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}

// 			// For error cases, check that the error message contains the expected substring
// 			if tt.wantErr && err != nil && !strings.Contains(err.Error(), tt.expectedErr) {
// 				t.Errorf("TransformJSONToSchemaFormat() error = %v, expected to contain %v", err, tt.expectedErr)
// 				return
// 			}

// 			// For success cases, validate the output format
// 			if !tt.wantErr {
// 				// Check that result is valid JSON
// 				var resultMap map[string]interface{}
// 				if err := json.Unmarshal([]byte(result), &resultMap); err != nil {
// 					t.Errorf("TransformJSONToSchemaFormat returned invalid JSON: %v", err)
// 					return
// 				}

// 				// Check that result has exactly two keys: schema and schemaType
// 				if len(resultMap) != 2 {
// 					t.Errorf("Expected result to have exactly 2 keys, got %d keys", len(resultMap))
// 				}

// 				// Check that schema key exists and is the original JSON string
// 				schema, ok := resultMap["schema"]
// 				if !ok {
// 					t.Errorf("Result is missing required 'schema' field")
// 				} else {
// 					// The schema field should be the original JSON as a string
// 					schemaStr, ok := schema.(string)
// 					if !ok {
// 						t.Errorf("'schema' field is not a string, got %T", schema)
// 					} else if schemaStr != tt.input {
// 						t.Errorf("Expected schema to be %q, got %q", tt.input, schemaStr)
// 					}
// 				}

// 				// Check that schemaType key exists and is "JSON"
// 				schemaType, ok := resultMap["schemaType"]
// 				if !ok {
// 					t.Errorf("Result is missing required 'schemaType' field")
// 				} else {
// 					// The schemaType field should be "JSON"
// 					schemaTypeStr, ok := schemaType.(string)
// 					if !ok {
// 						t.Errorf("'schemaType' field is not a string, got %T", schemaType)
// 					} else if schemaTypeStr != "JSON" {
// 						t.Errorf("Expected schemaType to be 'JSON', got %q", schemaTypeStr)
// 					}
// 				}

// 				// Also verify by creating the expected output manually and comparing
// 				expectedOutput := fmt.Sprintf(`{"schema":%q,"schemaType":"JSON"}`, tt.input)
// 				var expectedMap map[string]interface{}
// 				json.Unmarshal([]byte(expectedOutput), &expectedMap)

// 				expectedJSON, _ := json.Marshal(expectedMap)
// 				resultJSON, _ := json.Marshal(resultMap)

// 				if string(expectedJSON) != string(resultJSON) {
// 					t.Errorf("Expected output %s, got %s", expectedJSON, resultJSON)
// 				}
// 			}
// 		})
// 	}
// }

// // TestStructFieldsModification tests edge cases where the SchemaFormat struct might be modified
// func TestStructFieldsModification(t *testing.T) {
// 	// Define a mock struct that matches SchemaFormat but missing fields
// 	mockWithoutSchema := struct {
// 		SchemaType string `json:"schemaType"`
// 	}{
// 		SchemaType: "JSON",
// 	}

// 	mockWithoutSchemaType := struct {
// 		Schema string `json:"schema"`
// 	}{
// 		Schema: `{"test": true}`,
// 	}

// 	// This test verifies that our validation would catch missing fields
// 	// if they were hypothetically removed from the struct
// 	testCases := []struct {
// 		name        string
// 		mockStruct  interface{}
// 		expectedErr string
// 	}{
// 		{
// 			name:        "missing schema field",
// 			mockStruct:  mockWithoutSchema,
// 			expectedErr: "output missing required 'schema' field",
// 		},
// 		{
// 			name:        "missing schemaType field",
// 			mockStruct:  mockWithoutSchemaType,
// 			expectedErr: "output missing required 'schemaType' field",
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			// Marshal the mock struct
// 			jsonBytes, err := json.Marshal(tc.mockStruct)
// 			if err != nil {
// 				t.Fatalf("Failed to marshal mock struct: %v", err)
// 			}

// 			// Unmarshal to verify our struct would fail validation
// 			var outputMap map[string]interface{}
// 			if err := json.Unmarshal(jsonBytes, &outputMap); err != nil {
// 				t.Fatalf("Failed to unmarshal mock JSON: %v", err)
// 			}

// 			// Check for schema field
// 			_, hasSchema := outputMap["schema"]
// 			_, hasSchemaType := outputMap["schemaType"]

// 			// Simulate the validation from TransformJSONToSchemaFormat
// 			// At least one of these should fail based on our test cases
// 			if !hasSchema {
// 				if !strings.Contains(tc.expectedErr, "schema") {
// 					t.Errorf("Expected error about missing schema field, but got %q", tc.expectedErr)
// 				}
// 			}

// 			if !hasSchemaType {
// 				if !strings.Contains(tc.expectedErr, "schemaType") {
// 					t.Errorf("Expected error about missing schemaType field, but got %q", tc.expectedErr)
// 				}
// 			}

// 			// Reflect test to verify that SchemaFormat properly defines both fields
// 			schemaFormat := SchemaFormat{}
// 			schemaFormatType := reflect.TypeOf(schemaFormat)

// 			// Verify SchemaFormat has both required fields
// 			hasSchemaField := false
// 			hasSchemaTypeField := false
// 			for i := 0; i < schemaFormatType.NumField(); i++ {
// 				field := schemaFormatType.Field(i)
// 				jsonTag := field.Tag.Get("json")
// 				if jsonTag == "schema" {
// 					hasSchemaField = true
// 				}
// 				if jsonTag == "schemaType" {
// 					hasSchemaTypeField = true
// 				}
// 			}

// 			if !hasSchemaField {
// 				t.Error("SchemaFormat struct is missing 'schema' field with json tag")
// 			}

// 			if !hasSchemaTypeField {
// 				t.Error("SchemaFormat struct is missing 'schemaType' field with json tag")
// 			}
// 		})
// 	}
// }

// // TestFieldRequirements tests that the function properly validates required fields
// func TestFieldRequirements(t *testing.T) {
// 	// Create a helpers instance for testing but don't actually use it
// 	// since we're directly testing the validation logic
// 	_ = &Helpers{logger: slog.New(slog.NewTextHandler(io.Discard, nil))}
//
// 	// Test with a mock implementation that simulates missing keys
// 	testMissingKey := func(keyToOmit string) error {
// 		// Create a map that will be missing one key
// 		schemaMap := map[string]interface{}{
// 			"schema":     `{"test": true}`,
// 			"schemaType": "JSON",
// 		}

// 		// Remove the specified key
// 		delete(schemaMap, keyToOmit)

// 		// Convert to JSON
// 		mapJSON, _ := json.Marshal(schemaMap)

// 		// Unmarshal back to map for validation
// 		var resultMap map[string]interface{}
// 		if err := json.Unmarshal(mapJSON, &resultMap); err != nil {
// 			return err
// 		}

// 		// Use key, ok idiom to check for required fields
// 		if _, ok := resultMap["schema"]; !ok {
// 			return fmt.Errorf("missing required field: schema")
// 		}

// 		if _, ok := resultMap["schemaType"]; !ok {
// 			return fmt.Errorf("missing required field: schemaType")
// 		}

// 		return nil
// 	}

// 	// Test schema missing
// 	err := testMissingKey("schema")
// 	if err == nil {
// 		t.Error("Expected error for missing schema field, but got none")
// 	} else if !strings.Contains(err.Error(), "missing required field: schema") {
// 		t.Errorf("Expected error about missing schema field, got: %v", err)
// 	}

// 	// Test schemaType missing
// 	err = testMissingKey("schemaType")
// 	if err == nil {
// 		t.Error("Expected error for missing schemaType field, but got none")
// 	} else if !strings.Contains(err.Error(), "missing required field: schemaType") {
// 		t.Errorf("Expected error about missing schemaType field, got: %v", err)
// 	}

// 	// Test the actual function with real input
// 	validJSON := `{"test": true}`

// 	// Create function similar to TransformJSONToSchemaFormat but simulate missing keys
// 	makeTestFn := func(keyToOmit string) func(string) (string, error) {
// 		return func(jsonStr string) (string, error) {
// 			var jsonObj interface{}
// 			if err := json.Unmarshal([]byte(jsonStr), &jsonObj); err != nil {
// 				return "", fmt.Errorf("invalid JSON input: %v", err)
// 			}

// 			// Create schema map directly
// 			schemaMap := map[string]interface{}{
// 				"schema":     jsonStr,
// 				"schemaType": "JSON",
// 			}

// 			// Omit the specified key
// 			if keyToOmit != "" {
// 				delete(schemaMap, keyToOmit)
// 			}

// 			// Marshal to JSON
// 			tmpJSON, err := json.Marshal(schemaMap)
// 			if err != nil {
// 				return "", fmt.Errorf("error formatting schema: %v", err)
// 			}

// 			// Unmarshal back to map for validation with key, ok idiom
// 			var resultMap map[string]interface{}
// 			if err := json.Unmarshal(tmpJSON, &resultMap); err != nil {
// 				return "", fmt.Errorf("error validating schema format: %v", err)
// 			}

// 			// Use key, ok idiom to check for required fields
// 			if _, ok := resultMap["schema"]; !ok {
// 				return "", fmt.Errorf("missing required field: schema")
// 			}

// 			if _, ok := resultMap["schemaType"]; !ok {
// 				return "", fmt.Errorf("missing required field: schemaType")
// 			}

// 			return string(tmpJSON), nil
// 		}
// 	}

// 	// Test with schema field missing
// 	missingSchemaFn := makeTestFn("schema")
// 	_, err = missingSchemaFn(validJSON)
// 	if err == nil {
// 		t.Error("Expected error when schema field is missing")
// 	} else if !strings.Contains(err.Error(), "missing required field: schema") {
// 		t.Errorf("Expected error about missing schema field, got: %v", err)
// 	}

// 	// Test with schemaType field missing
// 	missingSchemaTypeFn := makeTestFn("schemaType")
// 	_, err = missingSchemaTypeFn(validJSON)
// 	if err == nil {
// 		t.Error("Expected error when schemaType field is missing")
// 	} else if !strings.Contains(err.Error(), "missing required field: schemaType") {
// 		t.Errorf("Expected error about missing schemaType field, got: %v", err)
// 	}
// }
