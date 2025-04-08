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
)

func handleHome(w http.ResponseWriter, r *http.Request) {
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

func handleSchema(w http.ResponseWriter, r *http.Request) {
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
	}{
		SubjectName: subjectName,
		Version:     version,
		SchemaID:    id,
	}
	t.Execute(w, data)
}

func handleTestSchemaPost(w http.ResponseWriter, r *http.Request) {

	subjectName := r.URL.Query().Get("topic")
	version := r.URL.Query().Get("version")
	id := r.URL.Query().Get("id")

	if subjectName == "" || version == "" || id == "" {
		log.Printf("Missing required parameters: subjectName=%s, version=%s, id=%s", subjectName, version, id)
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	// Convert version to integer
	versionInt, err := strconv.Atoi(version)
	if err != nil {
		log.Printf("Error converting version to integer: %v", err)
		http.Error(w, "Invalid version number", http.StatusBadRequest)
		return
	}

	log.Printf("Testing schema for subject: %s, version: %d, id: %s", subjectName, versionInt, id)
	log.Printf("Body: %s", string(body))

	// Test the schema
	isCompatible, statusCode, message, err := testSchema(subjectName, versionInt, string(body))
	if err != nil {
		log.Printf("Error testing schema: %v", err)
		http.Error(w, "Error testing schema", http.StatusInternalServerError)
		return
	}

	// Check if the message indicates compatibility couldn't be determined
	compatibilityDetermined := true
	if statusCode != http.StatusOK &&
		(statusCode == http.StatusInternalServerError ||
			strings.Contains(message, "error parsing") ||
			strings.Contains(message, "unexpected status")) {
		compatibilityDetermined = false
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

	// Prepare response
	response := map[string]interface{}{
		"is_compatible":            isCompatible,
		"compatibility_determined": compatibilityDetermined,
		"status_code":              statusCode,
		"error_code":               errorCode,
		"message":                  message,
	}

	// Set content type and encode response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("Schema test result: is_compatible=%t, compatibility_determined=%t, status_code=%d, error_code=%d, message=%s",
		isCompatible, compatibilityDetermined, statusCode, errorCode, message)
}

func handleTestSchemaAPI(w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if requestData.Subject == "" || requestData.Version == "" || requestData.Id == "" || requestData.JSON == "" {
		log.Printf("Missing required fields in API request")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Convert version to integer
	versionInt, err := strconv.Atoi(requestData.Version)
	if err != nil {
		log.Printf("Error converting version to integer: %v", err)
		http.Error(w, "Invalid version number", http.StatusBadRequest)
		return
	}

	log.Printf("Testing schema for subject: %s, version: %d, id: %s", requestData.Subject, versionInt, requestData.Id)

	// Test the schema
	isCompatible, statusCode, message, err := testSchema(requestData.Subject, versionInt, requestData.JSON)
	if err != nil {
		log.Printf("Error testing schema: %v", err)
		http.Error(w, "Error testing schema", http.StatusInternalServerError)
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

	// Prepare response using camelCase to match frontend expectations
	response := map[string]interface{}{
		"isCompatible": isCompatible,
		"httpStatus":   statusCode,
		"errorCode":    errorCode,
		"message":      message,
	}

	// Set content type and encode response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("API Schema test result: isCompatible=%t, httpStatus=%d, errorCode=%d, message=%s",
		isCompatible, statusCode, errorCode, message)
}
