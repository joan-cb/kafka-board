package confluentRegistryAPI

import (
	"bytes"
	"fmt"
	"net/http"
)

func createTestSchemaRequest(subjectName string, version int, testJSON string, baseRegistryURL string) (*http.Request, error) {
	requestURL := fmt.Sprintf("%s/compatibility/subjects/%s/versions/%d",
		baseRegistryURL, subjectName, version)

	req, err := http.NewRequest(
		"POST",
		requestURL,
		bytes.NewBuffer([]byte(testJSON)),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
