package handlers

import (
	"kafka-board/confluentRegistryAPI"
	"kafka-board/types"
	"log/slog"
)

type Handler struct {
	abstractRegistryAPI confluentRegistryAPI.RegistryAPI
	logger              *slog.Logger
}

// returnHandler creates and returns a new handler that implements registryAPICalls
// It can be extended to accept configuration options like base URLs, credentials, etc.
func ReturnHandler(logger *slog.Logger) *Handler {
	return &Handler{
		logger:              logger,
		abstractRegistryAPI: *confluentRegistryAPI.ReturnRegistryAPI(logger),
	}
}

type RegistryAPICalls interface {
	// API methods
	ReadeturnSubjects() ([]string, error)
	returnSubjectConfigs(subjectNames []string) ([]types.SubjectConfigInterface, error)
	GetGlobalConfig() (types.GlobalConfig, error)
	GetSchemas(subjectName string) ([]types.Schema, error)
	TestSchema(subjectName string, version int, testJSON string) (types.Response, error)
	GetSchema(id string) (types.Schema, error)
}
