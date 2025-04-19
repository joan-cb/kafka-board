package handlers

import (
	"encoding/json"
	"kafka-board/types"
	"net/http"
)

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
