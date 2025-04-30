package confluentRegistryAPI

import (
	"fmt"
	"net/http"
	"os"
)

func (r *RegistryAPI) deleteDefaultConfig() error {
	//to do -- > add delete default config implementation
	r.logger.Debug("deleteDefaultConfig - Deleting default config")
	return nil
}

func (r *RegistryAPI) deleteAllSubjects(subjectNames []string) (string, error) {
	client := &http.Client{}

	// Use the environment variable for Schema Registry URL
	registryURL := os.Getenv("SCHEMA_REGISTRY_URL")
	if registryURL == "" {
		registryURL = "http://schema-registry:8081"
	}

	baseURL := fmt.Sprintf("%s/subjects", registryURL)
	r.logger.Debug("DeleteAllSubjects - Using Schema Registry URL", "url", baseURL)

	for _, subjectName := range subjectNames {
		req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", baseURL, subjectName), nil)

		if err != nil {
			return "", err
		}
		resp, err := client.Do(req)

		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			r.logger.Error("DeleteAllSubjects - Unexpected status code",
				"status", resp.StatusCode,
				"subject", subjectName,
				"url", fmt.Sprintf("%s/%s", baseURL, subjectName))
			return fmt.Sprintf("%s/%s", baseURL, subjectName), fmt.Errorf("DeleteAllSubjects -unexpected status code: %d for subject: %s", resp.StatusCode, subjectName)
		}
		r.logger.Debug("DeleteAllSubjects - Subject deleted",
			"subject", subjectName)
	}

	return "", nil
}
