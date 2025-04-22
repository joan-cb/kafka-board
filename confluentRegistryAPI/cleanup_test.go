package confluentRegistryAPI

import (
	"fmt"
	"net/http"
)

func (r *RegistryAPI) deleteDefaultConfig() error {
	//to do -- > add delete default config implementation
	r.logger.Debug("deleteDefaultConfig - Deleting default config")
	return nil
}

func (r *RegistryAPI) deleteAllSubjects(subjectNames []string) (string, error) {
	client := &http.Client{}
	baseURL := "http://localhost:8090/subjects"

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
