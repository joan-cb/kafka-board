package handlers

// MockRegistryAPI mocks the RegistryAPI for testing
// type MockRegistryAPI struct {
// 	mockSchema types.Schema
// 	logger     *slog.Logger
// }

// func (m *MockRegistryAPI) ReturnSubjects() ([]string, error) {
// 	return []string{}, nil
// }

// func (m *MockRegistryAPI) ReturnSubjectConfigs(subjectNames []string) ([]types.SubjectConfigInterface, error) {
// 	return []types.SubjectConfigInterface{}, nil
// }

// func (m *MockRegistryAPI) GetGlobalConfig() (types.GlobalConfig, error) {
// 	return types.GlobalConfig{}, nil
// }

// func (m *MockRegistryAPI) GetSchemas(subjectName string) ([]types.Schema, error) {
// 	return []types.Schema{m.mockSchema}, nil
// }

// func (m *MockRegistryAPI) TestSchema(subjectName string, version int, testJSON string) (types.Response, error) {
// 	return types.Response{}, nil
// }

// func (m *MockRegistryAPI) GetSchema(id string) (types.Schema, error) {
// 	return m.mockSchema, nil
// }

// // TestHandler extends the real handler for testing
// type TestHandler struct {
// 	handler     *Handler
// 	mockAPI     *MockRegistryAPI
// 	logger      *slog.Logger
// 	mockHelpers *helpers.Helpers
// }
