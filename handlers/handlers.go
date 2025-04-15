package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"kafka-board/helpers"
	"kafka-board/types"
)

var falseVal = false
var trueVal = true

// Handle the home page load
func (h *Handler) HandleHomePage(w http.ResponseWriter, r *http.Request) {

	// First get all subjects
	subjects, err := h.abstractRegistryAPI.ReturnSubjects()

	if err != nil {
		h.logger.Debug("HandleHomePage - Error fetching subjects", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	//Fetch Global Config
	globalConfig, err := h.abstractRegistryAPI.GetGlobalConfig()

	if err != nil {
		h.logger.Debug("HandleHomePage - Error fetching global config",
			"error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	// Then get configs for all subjects
	configs, err := h.abstractRegistryAPI.ReturnSubjectConfigs(subjects)

	if err != nil {
		h.logger.Debug("HandleHomePage - Error fetching configs",
			"error", err)
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

	h.logger.Debug("HandleHomePage - Home page data",
		"data", data)

	t.Execute(w, data)
}

// Handle the schema page load
func (h *Handler) HandleSchemaPage(w http.ResponseWriter, r *http.Request) {
	subjectName := r.URL.Query().Get("topic")

	if subjectName == "" {
		h.logger.Debug("HandleSchemaPage - Subject name is required",
			"error", "subjectName is an empty string")

		return
	}

	schemas, err := h.abstractRegistryAPI.GetSchemas(subjectName)
	if err != nil {
		h.logger.Debug("HandleSchemaPage - Error fetching schemas", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	funcMap := template.FuncMap{
		"formatJSON": func(s string) string {
			var result interface{}
			if err := json.Unmarshal([]byte(s), &result); err != nil {
				h.logger.Debug("HandleSchemaPage - Error formatting JSON", "error", err)

				return s // Return original string if not valid JSON
			}
			formatted, err := json.MarshalIndent(result, "", "    ")
			if err != nil {
				h.logger.Error("HandleSchemaPage - Error formatting JSON", "error", err)

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

	h.logger.Debug("HandleSchemaPage - Schema data",
		"data", data)

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
		h.logger.Debug("HandleTestSchemaGet - Missing required parameters",
			"error", "subjectName, version, or id is an empty string")

		http.Error(w, "Missing required parameters", http.StatusBadRequest)

		return
	}

	// Get schemas for the subject
	schemas, err := h.abstractRegistryAPI.GetSchemas(subjectName)
	if err != nil {
		h.logger.Debug("HandleTestSchemaGet - Error fetching schemas",
			"error", err)
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
		h.logger.Debug("HandleTestSchemaGet - Schema not found",
			"error", "targetSchema.Version is 0 which is invalid")
		return
	}

	// Define the funcMap for the template to pretty-print the JSON schema
	funcMap := template.FuncMap{
		"formatJSON": func(s string) string {
			var result interface{}
			if err := json.Unmarshal([]byte(s), &result); err != nil {
				h.logger.Debug("HandleTestSchemaGet - Error formatting JSON",
					"error", err)
				return s
			}
			formatted, err := json.MarshalIndent(result, "", "    ")
			if err != nil {
				h.logger.Debug("HandleTestSchemaGet - Error formatting JSON",
					"error", err)
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
	h.logger.Debug("HandleTestSchemaGet - Schema data",
		"data", data)

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
		response := helpers.CreateResponseObject(
			nil,
			fmt.Sprintf("Error parsing JSON request: %v", err),
			http.StatusBadRequest,
			http.StatusBadRequest,
		)

		h.logger.Debug("HandleTestSchemaPost - Error parsing JSON request",
			"error", err)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Validate required fields
	if requestData.Subject == "" || requestData.Version == "" || requestData.Id == "" || requestData.JSON == "" {
		response := helpers.CreateResponseObject(
			nil,
			"Missing required fields",
			http.StatusBadRequest,
			http.StatusBadRequest,
		)

		h.logger.Debug("HandleTestSchemaPost - Missing required fields",
			"error", "requestData.Subject, requestData.Version, requestData.Id, or requestData.JSON is an empty string")

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Convert version to integer
	versionInt, err := strconv.Atoi(requestData.Version)
	if helpers.CheckErr(err) {
		response := helpers.CreateResponseObject(
			nil,
			fmt.Sprintf("Error parsing version number: %v", err),
			http.StatusBadRequest,
			http.StatusBadRequest,
		)

		h.logger.Debug("HandleTestSchemaPost - Error parsing version number",
			"error", err)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	h.logger.Debug("HandleTestSchemaPost - Testing schema",
		"subject", requestData.Subject,
		"version", versionInt,
		"json", requestData.JSON)

	// Test the schema
	resp, err := h.abstractRegistryAPI.TestSchema(requestData.Subject, versionInt, requestData.JSON)
	if helpers.CheckErr(err) {

		h.logger.Debug("HandleTestSchemaPost - Error testing schema",
			"error", err)

		helpers.SendJSONResponse(w, http.StatusInternalServerError, resp)

		return
	}

	// Ensure message has a value
	if resp.Message == "" {
		resp.Message = "None"
	}

	h.logger.Debug("HandleTestSchemaPost - Schema test successful",
		"message", resp.Message,
		"status", resp.StatusCode)

	helpers.SendJSONResponse(w, resp.StatusCode, resp)
}

// Called from the schema page to validate a payload against a schema
func (h *Handler) HandleValidatePayload(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Read and validate request body
	body, err := io.ReadAll(r.Body)
	if helpers.CheckErr(err) {
		response := helpers.CreateResponseObject(
			&falseVal,
			fmt.Sprintf("Error reading request body: %v", err),
			http.StatusBadRequest,
			0,
		)

		h.logger.Debug("HandleValidatePayload - Error reading request body",
			"error", err)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Parse the JSON request body
	var unmarshalledBody map[string]any
	err = json.Unmarshal(body, &unmarshalledBody)

	if helpers.CheckErr(err) {
		response := helpers.CreateResponseObject(
			&falseVal,
			fmt.Sprintf("Invalid JSON format in request body: %v", err),
			http.StatusBadRequest,
			0,
		)

		h.logger.Debug("HandleValidatePayload - Invalid JSON format in request body",
			"error", err)

		helpers.SendJSONResponse(w, http.StatusBadRequest, response)

		return
	}

	// Ensure payload key exists
	payloadRaw, ok := unmarshalledBody["payload"]
	if !ok {
		response := helpers.CreateResponseObject(
			&falseVal,
			"payload key expected in request body",
			http.StatusBadRequest,
			0,
		)

		h.logger.Debug("HandleValidatePayload - Payload key expected in request body",
			"error", "payload key expected in request body")

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
			response := helpers.CreateResponseObject(
				&falseVal,
				fmt.Sprintf("value of payload key is not valid JSON: %v", err),
				http.StatusBadRequest,
				0,
			)

			h.logger.Debug("HandleValidatePayload - Value of payload key is not valid JSON",
				"error", err)

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
		response := helpers.CreateResponseObject(
			&falseVal,
			fmt.Sprintf("Error retrieving schema: %v", err),
			http.StatusInternalServerError,
			0,
		)
		h.logger.Debug("HandleValidatePayload - Error retrieving schema",
			"error", err)
		helpers.SendJSONResponse(w, http.StatusInternalServerError, response)

		return
	}

	isValid, errors, err := helpers.ValidatePayload(payload, schema)

	if err != nil {
		h.logger.Debug("HandleValidatePayload - Error validating payload",
			"error", err)

		response := helpers.CreateResponseObject(
			&falseVal,
			fmt.Sprintf("Error validating payload: %v", err),
			http.StatusInternalServerError,
			0,
		)

		helpers.SendJSONResponse(w, http.StatusInternalServerError, response)

		return
	}

	// Create response based on validation result
	var response types.Response
	if !isValid {
		// Collect validation errors
		var errorMessages []string
		errorMessages = append(errorMessages, errors...)

		response = helpers.CreateResponseObject(
			&falseVal,
			strings.Join(errorMessages, "; "),
			http.StatusOK,
			0,
		)
	} else {
		response = helpers.CreateResponseObject(
			&trueVal,
			"Payload validates against schema",
			http.StatusOK,
			0,
		)
	}

	h.logger.Debug("HandleValidatePayload - Validation result",
		"valid", isValid,
		"errors", errors)

	helpers.SendJSONResponse(w, response.StatusCode, response)
}

func (h *Handler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	h.logger.Debug("HandleHealthCheck - Health check received")
	w.WriteHeader(http.StatusOK)
}
