package mocks

import "kafka-board/types"

type mockHandler struct {
	mockSchema types.Schema
}

func (m *mockHandler) getSchema(id string) (types.Schema, error) {
	return m.mockSchema, nil
}

func (m *mockHandler) returnSubjects() ([]string, error) {
	return []string{}, nil
}

func (m *mockHandler) returnSubjectConfigs(subjectNames []string) ([]types.SubjectConfigInterface, error) {
	return []types.SubjectConfigInterface{}, nil
}

func (m *mockHandler) getGlobalConfig() (types.GlobalConfig, error) {
	return types.GlobalConfig{}, nil
}

func (m *mockHandler) getSchemas(subjectName string) ([]types.Schema, error) {
	return []types.Schema{}, nil
}

func (m *mockHandler) testSchema(subjectName string, version int, testJSON string) (types.SchemaRegistryResponse, error) {
	return types.SchemaRegistryResponse{}, nil
}
