package helpers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Registry API configuration
var (
	// GetRegistryURL returns the base URL for the Schema Registry API
	// It reads from the SCHEMA_REGISTRY_URL environment variable or uses the default
	registryURL = func() string {
		if url := os.Getenv("SCHEMA_REGISTRY_URL"); url != "" {
			return url
		}
		return "http://schema-registry:8081"
	}()

	// Default HTTP client with reasonable timeouts
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          100,
			MaxConnsPerHost:       20,
			IdleConnTimeout:       30 * time.Second,
			DisableCompression:    false,
			ResponseHeaderTimeout: 5 * time.Second,
		},
	}
)

// CreateTestSchemaRequest creates a new HTTP request for testing schema compatibility
// against the Schema Registry API.
//
// Parameters:
//   - subjectName: The name of the subject (schema) to test compatibility against
//   - version: The version of the schema to test against
//   - testJSON: The JSON payload to test compatibility with
//
// Returns:
//   - *http.Request: A prepared HTTP request
//   - error: An error if request creation fails
func CreateTestSchemaRequest(subjectName string, version int, testJSON string) (*http.Request, error) {
	requestURL := fmt.Sprintf("%s/compatibility/subjects/%s/versions/%d",
		registryURL, subjectName, version)

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

// MakeHTTPRequest executes an HTTP request and returns the response
// This function uses a pre-configured HTTP client with timeouts
//
// Parameters:
//   - req: The HTTP request to execute
//
// Returns:
//   - *http.Response: The HTTP response
//   - error: An error if the request fails
func MakeHTTPRequest(req *http.Request) (*http.Response, error) {
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request to %s: %w", req.URL.String(), err)
	}

	return resp, nil
}

// MakeHTTPRequestWithContext executes an HTTP request with a context and returns the response
// This allows for request cancellation and timeout control at the call site
//
// Parameters:
//   - ctx: Context for the request execution
//   - req: The HTTP request to execute
//
// Returns:
//   - *http.Response: The HTTP response
//   - error: An error if the request fails or is canceled
func MakeHTTPRequestWithContext(ctx context.Context, req *http.Request) (*http.Response, error) {
	// Apply the context to the request
	req = req.WithContext(ctx)

	resp, err := httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.Canceled {
			return nil, fmt.Errorf("request was canceled: %w", err)
		}
		if ctx.Err() == context.DeadlineExceeded {
			return nil, fmt.Errorf("request timed out: %w", err)
		}
		return nil, fmt.Errorf("error making request to %s: %w", req.URL.String(), err)
	}

	return resp, nil
}

// ReadResponseBody reads and returns the response body
// It handles closing the response body automatically
//
// Parameters:
//   - resp: The HTTP response to read
//
// Returns:
//   - []byte: The response body as bytes
//   - error: An error if reading fails
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("cannot read body from nil response")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
