package main

type mockHandler struct {
	mockSchema Schema
}

func (m *mockHandler) getSchema(id string) (Schema, error) {
	return m.mockSchema, nil
}

func (m *mockHandler) returnSubjects() ([]string, error) {
	return []string{}, nil
}

func (m *mockHandler) returnSubjectConfigs(subjectNames []string) ([]SubjectConfigInterface, error) {
	return []SubjectConfigInterface{}, nil
}

func (m *mockHandler) getGlobalConfig() (GlobalConfig, error) {
	return GlobalConfig{}, nil
}

func (m *mockHandler) getSchemas(subjectName string) ([]Schema, error) {
	return []Schema{}, nil
}

func (m *mockHandler) testSchema(subjectName string, version int, testJSON string) (schemaRegistryResponse, error) {
	return schemaRegistryResponse{}, nil
}
