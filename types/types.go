package types

// Schema is the struct for the schema registry schema model
type Schema struct {
	Name       string `json:"name"`
	Subject    string `json:"subject"`
	Version    int    `json:"version"`
	Id         int    `json:"id"`
	SchemaType string `json:"schemaType"`
	Schema     string `json:"schema"`
}

type ConfigPayload struct {
	Compatibility string `json:"compatibility"`
}

// SubjectConfig is the struct for the subject config model
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

// SubjectGlobalConfig is the struct for subjects inheriting from the global config model
type SubjectGlobalConfig struct {
	Name               string `json:"name"`
	TakesGlobalDefault bool   `json:"takesGlobalDefault"`
}

// GlobalConfig is the struct for the global config model
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
// used
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

type Response struct {
	IsCompatible *bool  `json:"is_compatible"`
	ErrorCode    int    `json:"error_code"`
	Message      string `json:"message"`
	StatusCode   int    `json:"http_status"`
}

// SetDefaultNone sets "None" for any unpopulated string fields in the SubjectConfig
func (sc *SubjectConfig) SetDefaultNone() {
	if sc.Alias == "" {
		sc.Alias = "None"
	}
	if sc.CompatibilityLevel == "" {
		sc.CompatibilityLevel = "None set"
	}
	if sc.CompatibilityGroup == "" {
		sc.CompatibilityGroup = "None set"
	}
}

func (sc *GlobalConfig) SetDefaultNone() {
	if sc.Alias == "" {
		sc.Alias = "None"
	}
	if sc.CompatibilityLevel == "" {
		sc.CompatibilityLevel = "None set"
	}
	if sc.CompatibilityGroup == "" {
		sc.CompatibilityGroup = "None set"
	}
}
