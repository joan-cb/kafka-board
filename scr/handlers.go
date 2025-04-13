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
func (h *handler) handleHomePage(w http.ResponseWriter, r *http.Request) {
	// First get all subjects
	subjects, err := h.registryAPI.returnSubjects()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Fetch Global Config
	globalConfig, err := h.registryAPI.getGlobalConfig()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Then get configs for all subjects
	configs, err := h.registryAPI.returnSubjectConfigs(subjects)
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
func (h *handler) handleSchemaPage(w http.ResponseWriter, r *http.Request) {
	subjectName := r.URL.Query().Get("topic")
	if subjectName == "" {
		http.Error(w, "Subject name is required", http.StatusBadRequest)
		return
	}

	schemas, err := h.registryAPI.getSchemas(subjectName)
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
func (h *handler) handleTestSchema(w http.ResponseWriter, r *http.Request) {
	// Route to appropriate handler based on HTTP method
	switch r.Method {
	case http.MethodGet:
		h.handleTestSchemaGet(w, r)
	case http.MethodPost:
		h.handleTestSchemaPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle the test schema page load
func (h *handler) handleTestSchemaGet(w http.ResponseWriter, r *http.Request) {
	subjectName := r.URL.Query().Get("topic")
	version := r.URL.Query().Get("version")
	id := r.URL.Query().Get("id")

	if subjectName == "" || version == "" || id == "" {
		log.Printf("Missing required parameters: subjectName=%s, version=%s, id=%s", subjectName, version, id)
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Get schemas for the subject
	schemas, err := h.registryAPI.getSchemas(subjectName)
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
func (h *handler) handleTestSchemaPost(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestData struct {
		Subject string `json:"subject"`
		Version string `json:"version"`
		Id      string `json:"id"`
		JSON    string `json:"json"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if checkErr(err) {
		response := createSchemaRegistryResponse(
			nil,
			fmt.Sprintf("Error parsing JSON request: %v", err),
			http.StatusBadRequest,
			http.StatusBadRequest,
		)

		logger.Info("response sent to client",
			"function", "handleTestSchemaPost",
			"response", response)

		sendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Validate required fields
	if requestData.Subject == "" || requestData.Version == "" || requestData.Id == "" || requestData.JSON == "" {
		response := createSchemaRegistryResponse(
			nil,
			"Missing required fields",
			http.StatusBadRequest,
			http.StatusBadRequest,
		)

		logger.Info("response sent to client",
			"function", "handleTestSchemaPost",
			"response", response)
		sendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Convert version to integer
	versionInt, err := strconv.Atoi(requestData.Version)
	if checkErr(err) {
		response := createSchemaRegistryResponse(
			nil,
			fmt.Sprintf("Error parsing version number: %v", err),
			http.StatusBadRequest,
			http.StatusBadRequest,
		)
		logger.Info("response sent to client",
			"function", "handleTestSchemaPost",
			"response", response)

		sendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Test the schema
	resp, err := h.registryAPI.testSchema(requestData.Subject, versionInt, requestData.JSON)
	if checkErr(err) {
		logger.Info("response sent to client",
			"function", "handleTestSchemaPost",
			"response", resp)

		sendJSONResponse(w, http.StatusInternalServerError, resp)

		return
	}

	// Ensure message has a value
	if resp.Message == "" {
		resp.Message = "None"
	}
	logger.Info("response sent to client",
		"function", "handleTestSchemaPost",
		"response", resp)

	sendJSONResponse(w, resp.StatusCode, resp)
}

// Called from the schema page to validate a payload against a schema
func (h *handler) handleValidatePayload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Read and validate request body
	body, err := io.ReadAll(r.Body)
	if checkErr(err) {
		response := createPayloadResponse(
			false,
			fmt.Sprintf("Error reading request body: %v", err),
			http.StatusBadRequest,
		)
		logger.Info("response sent to client",
			"function", "handleValidatePayload - error reading request body",
			"response", response)

		sendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Parse the JSON request body
	var unmarshalledBody map[string]any
	err = json.Unmarshal(body, &unmarshalledBody)

	if checkErr(err) {
		response := createPayloadResponse(
			false,
			fmt.Sprintf("Invalid JSON format in request body: %v", err),
			http.StatusBadRequest,
		)
		logger.Info("response sent to client",
			"function", "handleValidatePayload - error unmarshalling request body",
			"response", response)

		sendJSONResponse(w, http.StatusBadRequest, response)

		return
	}
	logger.Debug("request body parsed",
		"function", "handleValidatePayload",
		"request body", unmarshalledBody)

	// Ensure payload key exists
	payloadRaw, ok := unmarshalledBody["payload"]
	if !ok {
		response := createPayloadResponse(
			false,
			"payload key expected in request body",
			http.StatusBadRequest,
		)

		logger.Info("response sent to client",
			"function", "handleValidatePayload - error payload key expected in request body",
			"response", response)

		sendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Process the payload based on its type
	var payload interface{}

	// Case 1: Payload is a string (need to parse as JSON)
	payloadStr, isString := payloadRaw.(string)
	if isString {
		err = json.Unmarshal([]byte(payloadStr), &payload)

		if checkErr(err) {
			response := createPayloadResponse(
				false,
				fmt.Sprintf("value of payload key is not valid JSON: %v", err),
				http.StatusBadRequest,
			)

			logger.Info("response sent to client",
				"function", "handleValidatePayload - value of payload key is not valid JSON",
				"response", response,
				"payload", payloadStr)

			sendJSONResponse(w, http.StatusBadRequest, response)

			return
		}
	} else {
		// Case 2: Payload is already a JSON object
		payload = payloadRaw
	}

	// Get the schema
	schema, err := h.registryAPI.getSchema(id)
	if checkErr(err) {
		response := createPayloadResponse(
			false,
			fmt.Sprintf("Error retrieving schema: %v", err),
			http.StatusInternalServerError,
		)

		logger.Info("response sent to client",
			"function", "handleValidatePayload ",
			"response", response)

		sendJSONResponse(w, http.StatusInternalServerError, response)

		return
	}

	// Create schema loader and validate
	schemaLoader := gojsonschema.NewStringLoader(schema.Schema)
	documentLoader := gojsonschema.NewGoLoader(payload)

	// Perform validation
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if checkErr(err) {
		response := createPayloadResponse(
			false,
			fmt.Sprintf("Error validating against schema: %v", err),
			http.StatusInternalServerError,
		)

		logger.Info("response sent to client",
			"function", "handleValidatePayload - error validating against schema",
			"response", response)

		sendJSONResponse(w, http.StatusInternalServerError, response)

		return
	}

	// Create response based on validation result
	var response payloadTestResponse
	if !result.Valid() {
		// Collect validation errors
		var errorMessages []string
		for _, err := range result.Errors() {
			errorMessages = append(errorMessages, err.String())
		}

		response = createPayloadResponse(
			false,
			strings.Join(errorMessages, "; "),
			http.StatusOK,
		)
	} else {
		response = createPayloadResponse(
			true,
			"Payload validates against schema",
			http.StatusOK,
		)
	}

	// Log and send the response
	logger.Info("response sent to client",
		"function", "handleValidatePayload - success",
		"response", response)

	sendJSONResponse(w, response.StatusCode, response)
}
