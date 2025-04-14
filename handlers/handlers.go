package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"kafka-board/helpers"
	"kafka-board/types"

	"github.com/xeipuuv/gojsonschema"
)

//TODO: Add STRUCTURED logger

// Handle the home page load
func (h *Handler) HandleHomePage(w http.ResponseWriter, r *http.Request) {
	// First get all subjects
	subjects, err := h.abstractRegistryAPI.ReturnSubjects()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Fetch Global Config
	globalConfig, err := h.abstractRegistryAPI.GetGlobalConfig()
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Then get configs for all subjects
	configs, err := h.abstractRegistryAPI.ReturnSubjectConfigs(subjects)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t := template.Must(template.New("home").Parse(homeTemplate))
	data := struct {
		Configs      []types.SubjectConfigInterface
		GlobalConfig types.GlobalConfig
	}{
		Configs:      configs,
		GlobalConfig: globalConfig,
	}
	t.Execute(w, data)
}

// Handle the schema page load
func (h *Handler) HandleSchemaPage(w http.ResponseWriter, r *http.Request) {
	subjectName := r.URL.Query().Get("topic")
	if subjectName == "" {
		http.Error(w, "Subject name is required", http.StatusBadRequest)
		return
	}

	schemas, err := h.abstractRegistryAPI.GetSchemas(subjectName)
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
		Schemas     []types.Schema
	}{
		SubjectName: subjectName,
		Schemas:     schemas,
	}
	t.Execute(w, data)
}

// Internal handler
func (h *Handler) HandleTestSchema(w http.ResponseWriter, r *http.Request) {
	// Route to appropriate handler based on HTTP method
	switch r.Method {
	case http.MethodGet:
		h.HandleTestSchemaGet(w, r)
	case http.MethodPost:
		h.HandleTestSchemaPost(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle the test schema page load
func (h *Handler) HandleTestSchemaGet(w http.ResponseWriter, r *http.Request) {
	subjectName := r.URL.Query().Get("topic")
	version := r.URL.Query().Get("version")
	id := r.URL.Query().Get("id")

	if subjectName == "" || version == "" || id == "" {
		log.Printf("Missing required parameters: subjectName=%s, version=%s, id=%s", subjectName, version, id)
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	// Get schemas for the subject
	schemas, err := h.abstractRegistryAPI.GetSchemas(subjectName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Find the specific schema version
	var targetSchema types.Schema
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
func (h *Handler) HandleTestSchemaPost(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var requestData struct {
		Subject string `json:"subject"`
		Version string `json:"version"`
		Id      string `json:"id"`
		JSON    string `json:"json"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if helpers.CheckErr(err) {
		response := helpers.CreateSchemaRegistryResponse(
			nil,
			fmt.Sprintf("Error parsing JSON request: %v", err),
			http.StatusBadRequest,
			http.StatusBadRequest,
		)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Validate required fields
	if requestData.Subject == "" || requestData.Version == "" || requestData.Id == "" || requestData.JSON == "" {
		response := helpers.CreateSchemaRegistryResponse(
			nil,
			"Missing required fields",
			http.StatusBadRequest,
			http.StatusBadRequest,
		)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Convert version to integer
	versionInt, err := strconv.Atoi(requestData.Version)
	if helpers.CheckErr(err) {
		response := helpers.CreateSchemaRegistryResponse(
			nil,
			fmt.Sprintf("Error parsing version number: %v", err),
			http.StatusBadRequest,
			http.StatusBadRequest,
		)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Test the schema
	resp, err := h.abstractRegistryAPI.TestSchema(requestData.Subject, versionInt, requestData.JSON)
	if helpers.CheckErr(err) {

		helpers.SendJSONResponse(w, http.StatusInternalServerError, resp)

		return
	}

	// Ensure message has a value
	if resp.Message == "" {
		resp.Message = "None"
	}

	helpers.SendJSONResponse(w, resp.StatusCode, resp)
}

// Called from the schema page to validate a payload against a schema
func (h *Handler) HandleValidatePayload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Read and validate request body
	body, err := io.ReadAll(r.Body)
	if helpers.CheckErr(err) {
		response := helpers.CreatePayloadResponse(
			false,
			fmt.Sprintf("Error reading request body: %v", err),
			http.StatusBadRequest,
		)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Parse the JSON request body
	var unmarshalledBody map[string]any
	err = json.Unmarshal(body, &unmarshalledBody)

	if helpers.CheckErr(err) {
		response := helpers.CreatePayloadResponse(
			false,
			fmt.Sprintf("Invalid JSON format in request body: %v", err),
			http.StatusBadRequest,
		)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Ensure payload key exists
	payloadRaw, ok := unmarshalledBody["payload"]
	if !ok {
		response := helpers.CreatePayloadResponse(
			false,
			"payload key expected in request body",
			http.StatusBadRequest,
		)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Process the payload based on its type
	var payload interface{}

	// Case 1: Payload is a string (need to parse as JSON)
	payloadStr, isString := payloadRaw.(string)
	if isString {
		err = json.Unmarshal([]byte(payloadStr), &payload)

		if helpers.CheckErr(err) {
			response := helpers.CreatePayloadResponse(
				false,
				fmt.Sprintf("value of payload key is not valid JSON: %v", err),
				http.StatusBadRequest,
			)

			helpers.SendJSONResponse(w, http.StatusBadRequest, response)

			return
		}
	} else {
		// Case 2: Payload is already a JSON object
		payload = payloadRaw
	}

	// Get the schema
	schema, err := h.abstractRegistryAPI.GetSchema(id)
	if helpers.CheckErr(err) {
		response := helpers.CreatePayloadResponse(
			false,
			fmt.Sprintf("Error retrieving schema: %v", err),
			http.StatusInternalServerError,
		)

		helpers.SendJSONResponse(w, http.StatusInternalServerError, response)

		return
	}

	// Create schema loader and validate
	schemaLoader := gojsonschema.NewStringLoader(schema.Schema)
	documentLoader := gojsonschema.NewGoLoader(payload)

	// Perform validation
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if helpers.CheckErr(err) {
		response := helpers.CreatePayloadResponse(
			false,
			fmt.Sprintf("Error validating against schema: %v", err),
			http.StatusInternalServerError,
		)

		helpers.SendJSONResponse(w, http.StatusInternalServerError, response)

		return
	}

	// Create response based on validation result
	var response types.PayloadTestResponse
	if !result.Valid() {
		// Collect validation errors
		var errorMessages []string
		for _, err := range result.Errors() {
			errorMessages = append(errorMessages, err.String())
		}

		response = helpers.CreatePayloadResponse(
			false,
			strings.Join(errorMessages, "; "),
			http.StatusOK,
		)
	} else {
		response = helpers.CreatePayloadResponse(
			true,
			"Payload validates against schema",
			http.StatusOK,
		)
	}

	helpers.SendJSONResponse(w, response.StatusCode, response)
}
