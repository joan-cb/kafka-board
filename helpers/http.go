package helpers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var baseRegistryURL = "http://schema-registry:8081"

// CreateTestSchemaRequest creates a new HTTP request for testing schema compatibility
func CreateTestSchemaRequest(subjectName string, version int, testJSON string) (*http.Request, error) {
	requestURL := baseRegistryURL + "/compatibility/subjects/" + subjectName + "/versions/" + strconv.Itoa(version)

	req, err := http.NewRequest(
		"POST",
		requestURL,
		bytes.NewBuffer([]byte(testJSON)),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.schemaregistry.v1+json")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// MakeHTTPRequest executes an HTTP request and returns the response
func MakeHTTPRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	return resp, nil
}

// ReadResponseBody reads and returns the response body
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	return body, nil
}
