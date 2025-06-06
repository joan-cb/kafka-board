package confluentRegistryAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"kafka-board/helpers"
	"kafka-board/types"
	"log/slog"
	"net/http"
	"os"
)

type RegistryAPI struct {
	logger          *slog.Logger
	baseRegistryURL string
}

func getBaseRegistryURL() string {
	if os.Getenv("REGISTRY_BASE_URL") == "" {
		return "http://localhost:8090"

	}
	return os.Getenv("REGISTRY_BASE_URL")
}

func ReturnRegistryAPI(logger *slog.Logger) *RegistryAPI {
	return &RegistryAPI{logger: logger, baseRegistryURL: getBaseRegistryURL()}
}

func (r *RegistryAPI) ReturnSubjects() ([]string, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects", r.baseRegistryURL), nil)
	if helpers.CheckErr(err) {
		r.logger.Debug("ReturnSubjects - Error creating request",
			"error", err)

		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")

	// Send request
	resp, err := client.Do(req)
	if helpers.CheckErr(err) {
		r.logger.Debug("ReturnSubjects - Error making request",
			"error", err)

		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		r.logger.Debug("ReturnSubjects - Unexpected status code",
			"status", resp.StatusCode)

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if helpers.CheckErr(err) {
		r.logger.Debug("ReturnSubjects - Error reading response",
			"error", err)

		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Parse JSON response
	var subjects []string
	if err := json.Unmarshal(body, &subjects); err != nil {
		r.logger.Debug("ReturnSubjects - Error parsing JSON",
			"body", string(body),
			"error", err)
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}
	r.logger.Debug("ReturnSubjects - Subjects returned by returnSubjects",
		"subjects", subjects)

	return subjects, nil
}

func (r *RegistryAPI) ReturnSubjectConfigs(subjectNames []string) ([]types.SubjectConfigInterface, error) {
	var configs []types.SubjectConfigInterface
	client := &http.Client{}

	for _, subjectName := range subjectNames {
		// Create request with URL-encoded subject name
		url := r.baseRegistryURL + "/config/" + subjectName
		r.logger.Debug("ReturnSubjectConfigs - Requesting URL",
			"url", url)

		req, err := http.NewRequest("GET", url, nil)

		if helpers.CheckErr(err) {
			r.logger.Debug("ReturnSubjectConfigs - Error creating request",
				"error", err)

			return nil, fmt.Errorf("error creating request: %v", err)
		}

		// Set headers
		req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")

		// Send request
		resp, err := client.Do(req)

		if helpers.CheckErr(err) {
			r.logger.Debug("ReturnSubjectConfigs - Error making request",
				"error", err)

			return nil, fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode == http.StatusNotFound {
			r.logger.Debug("ReturnSubjectConfigs - Subject config not found",
				"subject", subjectName)

			config := types.SubjectGlobalConfig{
				Name:               subjectName,
				TakesGlobalDefault: true,
			}

			configs = append(configs, config)

			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)

			r.logger.Debug("ReturnSubjectConfigs - Unexpected status code",
				"status", resp.StatusCode,
				"body", string(body))

			return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if helpers.CheckErr(err) {
			r.logger.Debug("ReturnSubjectConfigs - Error reading response",
				"error", err)

			return nil, fmt.Errorf("error reading response: %v", err)
		}

		// Parse JSON response
		config := types.SubjectConfig{
			Name: subjectName,
		}
		if helpers.CheckErr(json.Unmarshal(body, &config)) {
			r.logger.Debug("ReturnSubjectConfigs - Error parsing JSON",
				"error", err)

			return nil, fmt.Errorf("error parsing JSON: %v", err)
		}
		config.SetDefaultNone()
		r.logger.Debug("ReturnSubjectConfigs - Config returned by returnSubjectConfigs for subject",
			"subject", subjectName)

		configs = append(configs, config)
	}

	r.logger.Debug("ReturnSubjectConfigs - Configs returned by returnSubjectConfigs",
		"configs", configs)

	return configs, nil
}

func (r *RegistryAPI) GetGlobalConfig() (types.GlobalConfig, error) {
	client := &http.Client{}

	url := r.baseRegistryURL + "/config"
	req, err := http.NewRequest("GET", url, nil)

	if helpers.CheckErr(err) {
		r.logger.Debug("GetGlobalConfig -Error creating request",
			"error", err)

		return types.GlobalConfig{}, fmt.Errorf("error creating request: %v", err)
	}

	//Preparing Request
	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")
	resp, err := client.Do(req)

	if helpers.CheckErr(err) {
		r.logger.Debug("GetGlobalConfig - Error making request",
			"error", err)

		return types.GlobalConfig{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	//Checking Status Code
	if resp.StatusCode != http.StatusOK {
		r.logger.Debug("GetGlobalConfig - Unexpected status code",
			"status", resp.StatusCode)

		return types.GlobalConfig{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//Reading Response Body
	body, err := io.ReadAll(resp.Body)
	if helpers.CheckErr(err) {
		r.logger.Debug("GetGlobalConfig - Error reading response",
			"error", err)

		return types.GlobalConfig{}, fmt.Errorf("error reading response: %v", err)
	}

	//Parsing JSON Response
	globalConfig := types.GlobalConfig{
		Name: "Global Config",
	}

	//Setting Default Values
	globalConfig.SetDefaultNone()

	//Unmarshalling JSON Response
	if helpers.CheckErr(json.Unmarshal(body, &globalConfig)) {
		r.logger.Debug("GetGlobalConfig - Error parsing JSON",
			"error", err)

		return types.GlobalConfig{}, fmt.Errorf("error parsing JSON: %v", err)
	}
	r.logger.Debug("GetGlobalConfig - Global Config returned by getGlobalConfig",
		"config", globalConfig)

	return globalConfig, nil
}

func (r *RegistryAPI) GetSchemas(subjectName string) ([]types.Schema, error) {
	var allSchemas []types.Schema
	client := &http.Client{}

	url := r.baseRegistryURL + "/schemas"
	req, err := http.NewRequest("GET", url, nil)

	if helpers.CheckErr(err) {
		r.logger.Debug("GetSchemas - Error creating request",
			"error", err)

		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")

	resp, err := client.Do(req)

	if helpers.CheckErr(err) {
		r.logger.Debug("GetSchemas - Error making request",
			"error", err)

		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r.logger.Debug("GetSchemas - Unexpected status code",
			"status", resp.StatusCode)

		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if helpers.CheckErr(err) {
		r.logger.Debug("GetSchemas - Error reading response",
			"error", err)

		return nil, fmt.Errorf("error reading response: %v", err)
	}

	if err := json.Unmarshal(body, &allSchemas); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	// Filter schemas by subject
	var filteredSchemas []types.Schema
	for _, schema := range allSchemas {
		if schema.Subject == subjectName {
			filteredSchemas = append(filteredSchemas, schema)
		}
	}

	r.logger.Debug("GetSchemas - Filtered Schemas returned by getSchemas",
		"schemas", filteredSchemas)

	return filteredSchemas, nil
}

// transformJSONToSchemaFormat takes a JSON string and wraps it in the Schema Registry format
func (r *RegistryAPI) TestSchema(subjectName string, version int, testJSON string) (types.Response, error) {

	helper := helpers.ReturnHelpers(r.logger)
	payload, err := helper.TransformJSONToSchemaFormat(testJSON)

	if helpers.CheckErr(err) {
		r.logger.Debug("TestSchema - Error transforming JSON to Schema Registry format",
			"error", err)

		resp := helpers.CreateResponseObject(nil, fmt.Sprintf("Error transforming JSON to Schema Registry format. Invalid JSON string: %v", err), http.StatusBadRequest, http.StatusBadRequest)

		return resp, err
	}
	r.logger.Debug("TestSchema - Transformed JSON returned by transformJSONToSchemaFormat",
		"payload", payload)
	// Create the request
	req, err := createTestSchemaRequest(subjectName, version, payload, r.baseRegistryURL)

	if helpers.CheckErr(err) {
		r.logger.Debug("TestSchema - Error creating request",
			"error", err)

		resp := helpers.CreateResponseObject(nil, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError, http.StatusInternalServerError)
		return resp, err
	}

	// Make the request
	resp, err := helpers.MakeHTTPRequest(req)
	if helpers.CheckErr(err) {
		r.logger.Debug("TestSchema - Error making request",
			"error", err)

		resp := helpers.CreateResponseObject(nil, fmt.Sprintf("Error making request: %v", err), http.StatusInternalServerError, http.StatusInternalServerError)
		return resp, err
	}

	// Read the response body
	body, err := helpers.ReadResponseBody(resp)
	if helpers.CheckErr(err) {
		r.logger.Debug("TestSchema - Error reading response",
			"error", err)

		resp := helpers.CreateResponseObject(nil, fmt.Sprintf("Error reading response: %v", err), http.StatusInternalServerError, http.StatusInternalServerError)
		return resp, err
	}

	// Process the response
	result, err := helper.ProcessResponse(body, resp.StatusCode)
	if helpers.CheckErr(err) {
		r.logger.Debug("TestSchema - Error processing response",
			"error", err)

		resp := helpers.CreateResponseObject(nil, fmt.Sprintf("Error processing response: %v", err), http.StatusInternalServerError, http.StatusInternalServerError)
		return resp, err
	}

	// Ensure message is never empty
	if result.Message == "" {
		result.Message = "None"
	}

	// Ensure message is not too long as it will break the UI and Confluent returns long messages
	if len(result.Message) > 100 {
		result.Message = result.Message[:100] + "..."
	}

	// Handle nil IsCompatible pointer - this means the compatibility couldn't be determined
	if result.IsCompatible == nil {
		r.logger.Debug("TestSchema - IsCompatible is nil",
			"result", result)

		resp := helpers.CreateResponseObject(nil, result.Message, result.StatusCode, result.ErrorCode)
		return resp, nil
	}

	// Return the dereferenced bool value
	return result, nil
}

func (r *RegistryAPI) GetSchema(id string) (types.Schema, error) {
	schema := types.Schema{}
	client := &http.Client{}

	url := r.baseRegistryURL + "/schemas/ids/" + id
	req, err := http.NewRequest("GET", url, nil)
	if helpers.CheckErr(err) {
		r.logger.Debug("GetSchema - Error creating request",
			"error", err)

		return schema, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")

	resp, err := client.Do(req)
	if helpers.CheckErr(err) {
		r.logger.Debug("GetSchema - Error making request",
			"error", err)

		return schema, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		r.logger.Debug("GetSchema - Unexpected status code",
			"status", resp.StatusCode)

		return schema, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)

	if helpers.CheckErr(err) {
		r.logger.Debug("GetSchema - Error reading response",
			"error", err)

		return schema, fmt.Errorf("error reading response: %v", err)
	}

	if err := json.Unmarshal(body, &schema); err != nil {
		r.logger.Debug("GetSchema - Error parsing JSON",
			"error", err)

		return schema, fmt.Errorf("error parsing JSON: %v", err)
	}

	r.logger.Debug("GetSchema - Schema returned by getSchema",
		"schema", schema)

	return schema, nil
}
