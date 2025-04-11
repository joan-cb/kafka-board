package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/xeipuuv/gojsonschema"
)

// Handle the home page load
func handleHomePage(w http.ResponseWriter, r *http.Request) {
	// First get all subjects
	subjects, err := returnSubjects()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Fetch Global Config
	globalConfig, err := getGlobalConfig()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Then get configs for all subjects
	configs, err := returnSubjectConfigs(subjects)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := template.Must(template.New("home").Parse(homeTemplate))
	data := struct {
		Configs      []SubjectConfigInterface
		GlobalConfig GlobalConfig
	}{
		Configs:      configs,
		GlobalConfig: globalConfig,
	}
	t.Execute(w, data)
}

// Handle the schema page load
func handleSchemaPage(w http.ResponseWriter, r *http.Request) {
	subjectName := r.URL.Query().Get("topic")
	if subjectName == "" {
		http.Error(w, "Subject name is required", http.StatusBadRequest)
		return
	}

	schemas, err := getSchemas(subjectName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	funcMap := template.FuncMap{
		"formatJSON": func(s string) string {
			var result interface{}
			if err := json.Unmarshal([]byte(s), &result); err != nil {
				return s // Return original string if not valid JSON
			}
			formatted, err := json.MarshalIndent(result, "", "    ")
			if err != nil {
				return s // Return original string if formatting fails
			}
			return string(formatted)
		},
	}

	t := template.Must(template.New("schema").Funcs(funcMap).Parse(schemaTemplate))
	data := struct {
		SubjectName string
		Schemas     []Schema
	}{
		SubjectName: subjectName,
		Schemas:     schemas,
	}
	t.Execute(w, data)
}

// Internal handler
func handleTestSchema(w http.ResponseWriter, r *http.Request) {
	// Route to appropriate handler based on HTTP method
	switch r.Method {
	case http.MethodGet:
		handleTestSchemaGet(w, r)
	case http.MethodPost:
		handleTestSchemaPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle the test schema page load
func handleTestSchemaGet(w http.ResponseWriter, r *http.Request) {
	subjectName := r.URL.Query().Get("topic")
	version := r.URL.Query().Get("version")
	id := r.URL.Query().Get("id")

	if subjectName == "" || version == "" || id == "" {
		log.Printf("Missing required parameters: subjectName=%s, version=%s, id=%s", subjectName, version, id)
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Get schemas for the subject
	schemas, err := getSchemas(subjectName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find the specific schema version
	var targetSchema Schema
	for _, schema := range schemas {
		if fmt.Sprintf("%d", schema.Version) == version && fmt.Sprintf("%d", schema.Id) == id {
			targetSchema = schema
			break
		}
	}

	if targetSchema.Version == 0 {
		http.Error(w, "Schema not found", http.StatusNotFound)
		return
	}

	// Define the funcMap for the template to pretty-print the JSON schema
	funcMap := template.FuncMap{
		"formatJSON": func(s string) string {
			var result interface{}
			if err := json.Unmarshal([]byte(s), &result); err != nil {
				return s
			}
			formatted, err := json.MarshalIndent(result, "", "    ")
			if err != nil {
				return s
			}
			return string(formatted)
		},
	}

	t := template.Must(template.New("test").Funcs(funcMap).Parse(testSchemaTemplate))
	data := struct {
		SubjectName string
		Version     string
		SchemaID    string
		Schema      string
	}{
		SubjectName: subjectName,
		Version:     version,
		SchemaID:    id,
		Schema:      targetSchema.Schema,
	}
	t.Execute(w, data)
}

// Handle the test schema post request for testing the compatibility of a new schema against existing schema
func handleTestSchemaPost(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestData struct {
		Subject string `json:"subject"`
		Version string `json:"version"`
		Id      string `json:"id"`
		JSON    string `json:"json"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		log.Printf("Error parsing JSON request: %v", err)
		response := createSchemaRegistryResponse(
			false,
			fmt.Sprintf("Invalid JSON request: %v", err),
			http.StatusBadRequest,
			http.StatusBadRequest,
		)
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Validate required fields
	if requestData.Subject == "" || requestData.Version == "" || requestData.Id == "" || requestData.JSON == "" {
		log.Printf("Missing required fields in API request")
		response := createSchemaRegistryResponse(
			false,
			"Missing required fields",
			http.StatusBadRequest,
			http.StatusBadRequest,
		)
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// Convert version to integer
	versionInt, err := strconv.Atoi(requestData.Version)
	if err != nil {
		log.Printf("Error converting version to integer: %v", err)
		response := createSchemaRegistryResponse(
			false,
			fmt.Sprintf("Invalid version number: %v", err),
			http.StatusBadRequest,
			http.StatusBadRequest,
		)
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	log.Printf("Testing schema for subject: %s, version: %d, id: %s", requestData.Subject, versionInt, requestData.Id)

	// Test the schema
	isCompatible, statusCode, message, err := testSchema(requestData.Subject, versionInt, requestData.JSON)
	if err != nil {
		log.Printf("Error testing schema: %v", err)
		response := createSchemaRegistryResponse(
			false,
			fmt.Sprintf("Error testing schema: %v", err),
			http.StatusInternalServerError,
			http.StatusInternalServerError,
		)
		sendJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Determine error code based on status code
	errorCode := 0
	if statusCode >= 400 {
		errorCode = statusCode
	}

	// Ensure message has a value
	if message == "" {
		message = "None"
	}

	// Create response
	response := createSchemaRegistryResponse(
		isCompatible,
		message,
		statusCode,
		errorCode,
	)

	// Send JSON response
	sendJSONResponse(w, statusCode, response)

	log.Printf("API Schema test result: isCompatible=%t, httpStatus=%d, errorCode=%d, message=%s",
		isCompatible, statusCode, errorCode, message)
}

// Called from the schema page to validate a payload against a schema
func handleValidatePayload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Read and validate request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		response := createPayloadResponse(
			false,
			fmt.Sprintf("Error reading request body: %v", err),
			http.StatusBadRequest,
		)
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	var unmarshalledBody map[string]any
	if err := json.Unmarshal(body, &unmarshalledBody); err != nil {
		log.Printf("Invalid JSON in request body: %v", err)
		response := createPayloadResponse(
			false,
			fmt.Sprintf("Invalid JSON format in request body: %v", err),
			http.StatusBadRequest,
		)
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	payloadRaw, ok := unmarshalledBody["payload"]
	if !ok {
		log.Print("no payload key provided")
		response := createPayloadResponse(
			false,
			"payload key expected in request body",
			http.StatusBadRequest,
		)
		sendJSONResponse(w, http.StatusBadRequest, response)
		return
	}

	// If payload is a string, try to parse it as JSON
	var payload interface{}
	if payloadStr, isString := payloadRaw.(string); isString {
		if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
			log.Printf("Invalid JSON in payload string: %v", err)
			response := createPayloadResponse(
				false,
				fmt.Sprintf("Invalid JSON in payload string: %v", err),
				http.StatusBadRequest,
			)
			sendJSONResponse(w, http.StatusBadRequest, response)
			return
		}
	} else {
		// If it's already a JSON object, use it directly
		payload = payloadRaw
	}

	// Get the schema
	schema, err := getSchema(id)
	if err != nil {
		log.Printf("Error retrieving schema: %v", err)
		response := createPayloadResponse(
			false,
			fmt.Sprintf("Error retrieving schema: %v", err),
			http.StatusInternalServerError,
		)
		sendJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Create schema loader
	schemaLoader := gojsonschema.NewStringLoader(schema.Schema)
	documentLoader := gojsonschema.NewGoLoader(payload)

	// Perform validation
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		log.Printf("Error during schema validation: %v", err)
		response := createPayloadResponse(
			false,
			fmt.Sprintf("Error validating against schema: %v", err),
			http.StatusInternalServerError,
		)
		sendJSONResponse(w, http.StatusInternalServerError, response)
		return
	}

	// Check validation result
	if !result.Valid() {
		// Collect validation errors
		var errorMessages []string
		for _, err := range result.Errors() {
			errorMessages = append(errorMessages, err.String())
		}

		response := createPayloadResponse(
			false,
			strings.Join(errorMessages, "; "),
			http.StatusOK,
		)
		log.Print(response)
		sendJSONResponse(w, http.StatusOK, response)
		return
	}

	// If we get here, validation passed
	response := createPayloadResponse(
		true,
		"Payload validates against schema",
		http.StatusOK,
	)
	log.Print(response)
	sendJSONResponse(w, http.StatusOK, response)
}
