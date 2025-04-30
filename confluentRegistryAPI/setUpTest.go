package confluentRegistryAPI

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"kafka-board/helpers"
	"kafka-board/types"
	"net/http"
)

type testSubject struct {
	schemaStr          string
	takesDefaultConfig bool
	subjectName        string
	config             types.ConfigPayload
}

// var defaultConfig = `{"compatibility": "BACKWARD"}`

var newSubjects = []testSubject{
	{
		schemaStr:          `{"schema": "{\"type\": \"record\", \"name\": \"test\", \"fields\": [{\"name\": \"field1\", \"type\": \"string\"}]}"}`,
		takesDefaultConfig: false,
		subjectName:        "test-forward_transitive",
		config: types.ConfigPayload{
			Compatibility: "forward_transitive",
		},
	},
	{
		schemaStr:          `{"schema": "{\"type\": \"record\", \"name\": \"test\", \"fields\": [{\"name\": \"field_2\", \"type\": \"boolean\"}]}"}`,
		takesDefaultConfig: false,
		subjectName:        "test-none",
		config: types.ConfigPayload{
			Compatibility: "none",
		},
	},
	{
		schemaStr:          `{"schema": "{\"type\": \"record\", \"name\": \"test\", \"fields\": [{\"name\": \"field_2\", \"type\": \"boolean\"}]}"}`,
		takesDefaultConfig: false,
		subjectName:        "test-subject-global config",
		config:             types.ConfigPayload{},
	},
}

func (r *RegistryAPI) createTestSubject(testSubject testSubject) error {
	client := &http.Client{}
	requestURL := fmt.Sprintf("%s/subjects/%s/versions", r.baseRegistryURL, testSubject.subjectName)
	helper := helpers.ReturnHelpers(r.logger)

	r.logger.Debug("CreateTestSubject - Using Schema Registry URL",
		"url", requestURL)

	// Create helper instance first
	transformedSchema, err := helper.TransformJSONToSchemaFormat(testSubject.schemaStr)
	if err != nil {
		r.logger.Debug("CreateTestSubject - Error transforming schema",
			"error", err)
		return err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(transformedSchema)))
	if err != nil {
		r.logger.Debug("CreateTestSubject - Error creating request",
			"error", err)
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		r.logger.Debug("CreateTestSubject - Error making request",
			"error", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		r.logger.Debug("CreateTestSubject - Error creating subject",
			"status", resp.Status,
			"body", string(bodyBytes))
		return fmt.Errorf("failed to create subject: %s", resp.Status)
	}
	r.logger.Debug("CreateTestSubject - Response",
		"status", resp.Status,
		"body", resp.Body)

	return nil
}

func (r *RegistryAPI) createConfig(testSubject testSubject) error {
	if testSubject.takesDefaultConfig {
		return nil
	}
	client := &http.Client{}
	requestURL := fmt.Sprintf("%s/config/%s", r.baseRegistryURL, testSubject.subjectName)
	r.logger.Debug("CreateConfig - Using Schema Registry URL", "url", requestURL)

	requestPayload, err := json.Marshal(testSubject.config)
	if err != nil {
		r.logger.Debug("CreateConfig - Error marshalling config",
			"error", err)
		return err
	}

	req, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(requestPayload))
	if err != nil {
		r.logger.Debug("CreateConfig - Error creating request",
			"error", err)
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
