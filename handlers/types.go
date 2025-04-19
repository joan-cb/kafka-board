package handlers

import (
	"kafka-board/helpers"
	"kafka-board/types"
	"log/slog"
)

type handler struct {
	registryAPI registryAPICalls
	logger      *slog.Logger
	helpers     *helpers.Helpers
}

// returnHandler creates and returns a new handler that implements registryAPICalls
// It can be extended to accept configuration options like base URLs, credentials, etc.
func ReturnHandler(logger *slog.Logger, registryConcreteImplementation registryAPICalls) *handler {
	return &handler{
		logger:      logger,
		registryAPI: registryConcreteImplementation,
		helpers:     helpers.ReturnHelpers(logger),
	}
}

type registryAPICalls interface {
	// API methods
	ReturnSubjects() ([]string, error)
	ReturnSubjectConfigs(subjectNames []string) ([]types.SubjectConfigInterface, error)
	GetGlobalConfig() (types.GlobalConfig, error)
	GetSchemas(subjectName string) ([]types.Schema, error)
	TestSchema(subjectName string, version int, testJSON string) (types.Response, error)
	GetSchema(id string) (types.Schema, error)
}
