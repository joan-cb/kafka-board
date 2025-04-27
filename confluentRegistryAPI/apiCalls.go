package confluentRegistryAPI

import (
	"encoding/json"
	"fmt"
	"io"
	"kafka-board/helpers"
	"kafka-board/types"
	"log"
	"log/slog"
	"net/http"
	"os"
)

var baseRegistryURL = os.Getenv("REGISTRY_BASE_URL")

type RegistryAPI struct {
	logger *slog.Logger
}

func ReturnRegistryAPI(logger *slog.Logger) *RegistryAPI {
	return &RegistryAPI{logger: logger}
}

func (r *RegistryAPI) ReturnSubjects() ([]string, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects", baseRegistryURL), nil)
	if err != nil {
		r.logger.Debug("ReturnSubjects - Error creating request",
			"error", err)

		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
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
	if err != nil {
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
		url := baseRegistryURL + "/config/" + subjectName
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
	url := baseRegistryURL + "/config"
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
	url := baseRegistryURL + "/schemas"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
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
	return filteredSchemas, nil
}

// transformJSONToSchemaFormat takes a JSON string and wraps it in the Schema Registry format
func (r *RegistryAPI) TestSchema(subjectName string, version int, testJSON string) (types.Response, error) {

	//validate JSON and transform to Schema Registry format
	// Transform JSON to Schema Registry format
	helper := helpers.ReturnHelpers(r.logger)
	payload, err := helper.TransformJSONToSchemaFormat(testJSON)
	if err != nil {
		log.Printf("Error transforming JSON to Schema Registry format: %v", err)
		resp := helpers.CreateResponseObject(nil, fmt.Sprintf("Error transforming JSON to Schema Registry format. Invalid JSON string: %v", err), http.StatusBadRequest, http.StatusBadRequest)
		return resp, err
	}
	log.Printf("Transformed JSON returned by transformJSONToSchemaFormat: %s", payload)
	// Create the request
	req, err := helpers.CreateTestSchemaRequest(subjectName, version, payload)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		resp := helpers.CreateResponseObject(nil, fmt.Sprintf("Error creating request: %v", err), http.StatusInternalServerError, http.StatusInternalServerError)
		return resp, err
	}

	// Make the request
	resp, err := helpers.MakeHTTPRequest(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		resp := helpers.CreateResponseObject(nil, fmt.Sprintf("Error making request: %v", err), http.StatusInternalServerError, http.StatusInternalServerError)
		return resp, err
	}

	// Read the response body
	body, err := helpers.ReadResponseBody(resp)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		resp := helpers.CreateResponseObject(nil, fmt.Sprintf("Error reading response: %v", err), http.StatusInternalServerError, http.StatusInternalServerError)
		return resp, err
	}

	// Process the response
	result, err := helper.ProcessResponse(body, resp.StatusCode)
	if err != nil {
		log.Printf("Error processing response: %v", err)
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
		resp := helpers.CreateResponseObject(nil, result.Message, result.StatusCode, result.ErrorCode)
		return resp, nil
	}

	// Return the dereferenced bool value
	return result, nil
}

func (r *RegistryAPI) GetSchema(id string) (types.Schema, error) {

	schema := types.Schema{}

	client := &http.Client{}
	url := baseRegistryURL + "/schemas/ids/" + id
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return schema, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")

	resp, err := client.Do(req)
	if err != nil {
		return schema, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return schema, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return schema, fmt.Errorf("error reading response: %v", err)
	}

	if err := json.Unmarshal(body, &schema); err != nil {
		return schema, fmt.Errorf("error parsing JSON: %v", err)
	}

	log.Printf("Schema returned by getSchema: %v", schema)
	return schema, nil
}
