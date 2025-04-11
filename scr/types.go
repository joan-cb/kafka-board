package main

type Schema struct {
	Name       string `json:"name"`
	Subject    string `json:"subject"`
	Version    int    `json:"version"`
	Id         int    `json:"id"`
	SchemaType string `json:"schemaType"`
	Schema     string `json:"schema"`
}

type SubjectConfig struct {
	Name               string `json:"name"`
	Normalize          *bool  `json:"normalize"`
	Alias              string `json:"aliases"`
	CompatibilityLevel string `json:"compatibilityLevel"`
	CompatibilityGroup string `json:"compatibilityGroup"`
	DefaultMetadata    any    `json:"defaultMetadata"`
	OverrideMetadata   any    `json:"overrideMetadata"`
	DefaultRuleSet     any    `json:"defaultRuleSet"`
	OverrideRuleSet    any    `json:"overrideRuleSet"`
}

type SubjectGlobalConfig struct {
	Name               string `json:"name"`
	TakesGlobalDefault bool   `json:"takesGlobalDefault"`
}

// to do: handle values of type any in the UI
type GlobalConfig struct {
	Name               string `json:"name"`
	Normalize          *bool  `json:"normalize"`
	Alias              string `json:"aliases"`
	CompatibilityLevel string `json:"compatibilityLevel"`
	CompatibilityGroup string `json:"compatibilityGroup"`
	DefaultMetadata    any    `json:"defaultMetadata"`
	OverrideMetadata   any    `json:"overrideMetadata"`
	DefaultRuleSet     any    `json:"defaultRuleSet"`
	OverrideRuleSet    any    `json:"overrideRuleSet"`
}

// SubjectConfigInterface defines the common interface for both config types
type SubjectConfigInterface interface {
	GetName() string
}

// Implement the interface for SubjectConfig
func (sc SubjectConfig) GetName() string {
	return sc.Name
}

// Implement the interface for SubjectGlobalConfig
func (sgc SubjectGlobalConfig) GetName() string {
	return sgc.Name
}

type schemaRegistryResponse struct {
	IsCompatible *bool  `json:"is_compatible"`
	ErrorCode    int    `json:"error_code"`
	Message      string `json:"message"`
	HttpStatus   int    `json:"http_status"`
}

type payloadTestResponse struct {
	IsCompatible *bool  `json:"is_compatible"`
	Message      string `json:"message"`
	HttpStatus   int    `json:"http_status"`
}
