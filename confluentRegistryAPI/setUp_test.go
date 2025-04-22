package confluentRegistryAPI

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type testSubject struct {
	schemaStr          string
	takesDefaultConfig bool
	subjectName        string
}

var defaultConfig = `{"compatibility": "BACKWARD"}`

var newSubjects = []testSubject{
	{
		schemaStr:          `{"schema": "{\"type\": \"record\", \"name\": \"test\", \"fields\": [{\"name\": \"field1\", \"type\": \"string\"}]}"}`,
		takesDefaultConfig: true,
		subjectName:        "test-subject",
	},
	{
		schemaStr:          `{"schema": "{\"type\": \"record\", \"name\": \"test\", \"fields\": [{\"name\": \"field_2\", \"type\": \"boolean\"}]}"}`,
		takesDefaultConfig: false,
		subjectName:        "test-subject-2",
	},
}

func (r *RegistryAPI) createTestSubject(newSubjects []testSubject) error {
	client := &http.Client{}
	baseURL := "http://localhost:8090/subjects"

	// Simple Avro schema for testing
	for _, subjectName := range newSubjects {
		// Create request with schema in body
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s/versions", baseURL, subjectName.subjectName),
			strings.NewReader(subjectName.schemaStr))

		if err != nil {
			r.logger.Debug("CreateTestSubject - Error creating request",
				"error", err)
			return err
		}

		// Set content type for JSON payload
		req.Header.Set("Content-Type", "application/json")

		// Send request
		resp, err := client.Do(req)

		if err != nil {
			r.logger.Debug("CreateTestSubject - Error making request",
				"error", err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			r.logger.Debug("CreateTestSubject - Unexpected status code",
				"status", resp.StatusCode,
				"subject", subjectName,
				"response", string(body))
			r.logger.Error("CreateTestSubject - Unexpected status code",
				"status", resp.StatusCode,
				"subject", subjectName,
				"response", string(body))
			return fmt.Errorf("CreateTestSubject - unexpected status code: %d for subject: %s", resp.StatusCode, subjectName.subjectName)
		}

		r.logger.Debug("CreateTestSubject - Subject created",
			"subject", subjectName.subjectName)
	}

	return nil
}

func (r *RegistryAPI) createDefaultConfig(defaultConfig string) error {
	//to do
	return nil
}
