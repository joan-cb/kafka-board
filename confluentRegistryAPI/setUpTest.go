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
	success            bool
	takesDefaultConfig bool
	subjectName        string
	config             types.ConfigPayload
}

// var defaultConfig = `{"compatibility": "BACKWARD"}`

var newSubjects = []testSubject{
	{
		schemaStr: `{
  "schema": {
    "type": "object",
    "properties": {
      "field1": {
        "type": "string"
      },
      "field2": {
        "type": "integer"
      },
      "field3": {
        "type": "boolean"
      }
    },
    "required": ["field1"]
  }
}`,
		takesDefaultConfig: false,
		subjectName:        "test-forward",
		config: types.ConfigPayload{
			Compatibility: "forward",
		},
	},
	{
		schemaStr: `{
	"schema": {
		"type": "object",
		"properties": {
		"field1": {
			"type": "string"
		},
		"field2": {
			"type": "integer"
		},
		"field3": {
			"type": "boolean"
		}
		},
		"required": ["field1","field2"]
	}
	}`,
		takesDefaultConfig: false,
		subjectName:        "test-backward",
		config: types.ConfigPayload{
			Compatibility: "backward",
		},
	},
	{
		schemaStr: `{
	"schema": {
		"type": "object",
		"properties": {
		"field1": {
			"type": "string"
		},
		"field2": {
			"type": "integer"
		},
		"field3": {
			"type": "boolean"
		}
		},
		"required": ["field1","field2"]
	}
	}`,
		takesDefaultConfig: false,
		subjectName:        "test-full",
		config: types.ConfigPayload{
			Compatibility: "full",
		},
	},
}

func (r *RegistryAPI) createTestSubject(testSubject testSubject) error {

	// Create helper instance to transform the input string into the schema registry format
	helper := helpers.ReturnHelpers(r.logger)
	transformedSchema, err := helper.TransformJSONToSchemaFormat(testSubject.schemaStr)

	if err != nil {
		r.logger.Debug("CreateTestSubject - Error transforming schema",
			"error", err)

		return err
	}

	client := &http.Client{}
	requestURL := fmt.Sprintf("%s/subjects/%s/versions", r.baseRegistryURL, testSubject.subjectName)

	r.logger.Debug("CreateTestSubject - Using Schema Registry URL",
		"url", requestURL)

	// Create the request
	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer([]byte(transformedSchema)))

	if err != nil {
		r.logger.Debug("CreateTestSubject - Error creating request",
			"error", err)

		return err
	}
	req.Header.Set("Content-Type", "application/json")

	// Make the request
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

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create config: %s", resp.Status)
	}

	return nil
}
