package handlers

import (
	"fmt"
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
