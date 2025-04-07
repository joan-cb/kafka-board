package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var baseRegistryURL = "http://schema-registry:8081"

// setDefaultNone sets "None" for any unpopulated string fields in the SubjectConfig
func (sc *SubjectConfig) setDefaultNone() {
	if sc.Name == "" {
		sc.Name = "None"
	}
	if sc.Alias == "" {
		sc.Alias = "None"
	}
	if sc.CompatibilityLevel == "" {
		sc.CompatibilityLevel = "None"
	}
}

func (sc *GlobalConfig) setDefaultNone() {
	if sc.Alias == "" {
		sc.Alias = "None"
	}
	if sc.CompatibilityLevel == "" {
		sc.CompatibilityLevel = "None"
	}
	if sc.CompatibilityGroup == "" {
		sc.CompatibilityGroup = "None"
	}
}

func returnSubjects() ([]string, error) {
	// Create HTTP client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/subjects", baseRegistryURL), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	// Parse JSON response
	var subjects []string
	if err := json.Unmarshal(body, &subjects); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}
	log.Printf("Subjects returned by returnSubjects: %v", subjects)
	return subjects, nil
}

func returnSubjectConfigs(subjectNames []string) ([]SubjectConfigInterface, error) {
	var configs []SubjectConfigInterface
	client := &http.Client{}

	for _, subjectName := range subjectNames {
		// Create request with URL-encoded subject name
		url := baseRegistryURL + "/config/" + subjectName
		log.Printf("Requesting URL: %s", url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Printf("Error creating request: %v", err)
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		// Set headers
		req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("Error making request: %v", err)
			return nil, fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		// Check status code
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("Subject config not found: %s", subjectName)
			config := SubjectGlobalConfig{
				Name:               subjectName,
				TakesGlobalDefault: true,
			}
			configs = append(configs, config)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Printf("Unexpected status code: %d, body: %s", resp.StatusCode, string(body))
			return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response: %v", err)
			return nil, fmt.Errorf("error reading response: %v", err)
		}

		// Parse JSON response
		config := SubjectConfig{
			Name: subjectName,
		}
		if err := json.Unmarshal(body, &config); err != nil {
			return nil, fmt.Errorf("error parsing JSON: %v", err)
		}
		config.setDefaultNone()
		log.Printf("Config returned by returnSubjectConfigs for subject: %s", subjectName)
		configs = append(configs, config)
	}
	log.Printf("Configs returned by returnSubjectConfigs: %v", configs)
	return configs, nil
}

func getGlobalConfig() (GlobalConfig, error) {
	client := &http.Client{}
	url := baseRegistryURL + "/config"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return GlobalConfig{}, fmt.Errorf("error creating request: %v", err)
	}

	//Preparing Request
	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return GlobalConfig{}, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	//Checking Status Code
	if resp.StatusCode != http.StatusOK {
		return GlobalConfig{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//Reading Response Body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return GlobalConfig{}, fmt.Errorf("error reading response: %v", err)
	}

	//Parsing JSON Response
	config := GlobalConfig{
		Name: "Global Config",
	}

	//Setting Default Values
	config.setDefaultNone()

	//Unmarshalling JSON Response
	if err := json.Unmarshal(body, &config); err != nil {
		log.Println(err)
		return GlobalConfig{}, fmt.Errorf("error parsing JSON: %v", err)
	}
	log.Printf("Global Config returned by getGlobalConfig: %v", config)
	return config, nil
}

func getSchema(subjectName string) ([]Schema, error) {
	var allSchemas []Schema
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
	var filteredSchemas []Schema
	for _, schema := range allSchemas {
		if schema.Subject == subjectName {
			filteredSchemas = append(filteredSchemas, schema)
		}
	}
	log.Printf("Schemas returned by getSchema: %v", filteredSchemas)
	return filteredSchemas, nil
}

// transformJSONToSchemaFormat takes a JSON string and wraps it in the Schema Registry format

func testSchema(subjectName string, version int, testJSON string) (bool, int, string, error) {
	// Create the request
	req, err := createTestSchemaRequest(subjectName, version, testJSON)
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return false, http.StatusInternalServerError, fmt.Sprintf("Error creating request: %v", err), err
	}

	// Make the request
	resp, err := makeHTTPRequest(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return false, http.StatusInternalServerError, fmt.Sprintf("Error making request: %v", err), err
	}

	// Read the response body
	body, err := readResponseBody(resp)
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return false, http.StatusInternalServerError, fmt.Sprintf("Error reading response: %v", err), err
	}

	// Process the response
	response, err := processResponse(body, resp.StatusCode)
	if err != nil {
		log.Printf("Error processing response: %v", err)
		return false, http.StatusInternalServerError, fmt.Sprintf("Error processing response: %v", err), err
	}

	// Ensure message is never empty
	if response.Message == "" {
		response.Message = "None"
	}

	// Ensure message is not too long as it will break the UI and Confluent returns long messages
	if len(response.Message) > 100 {
		response.Message = response.Message[:100] + "..."
	}

	// Handle nil IsCompatible pointer - this means the compatibility couldn't be determined
	if response.IsCompatible == nil {
		return false, response.HttpStatus, response.Message, nil
	}

	// Return the dereferenced bool value
	return *response.IsCompatible, response.HttpStatus, response.Message, nil
}
