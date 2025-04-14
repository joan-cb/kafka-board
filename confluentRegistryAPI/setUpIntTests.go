package confluentRegistryAPI

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"
// )

// // Setup a test subject with schema for testing
// func setupTestSubject(api *registryAPI, name string, schema string) (int, error) {
// 	// This would require implementing a registration method
// 	// Example implementation might call the Schema Registry REST API directly
// 	client := &http.Client{}

// 	url := baseRegistryURL + "/subjects/" + name + "/versions"
// 	payload := fmt.Sprintf(`{"schema": %q}`, schema)

// 	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
// 	if err != nil {
// 		return 0, err
// 	}

// 	req.Header.Set("Content-Type", "application/json")
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return 0, fmt.Errorf("unexpected status: %d", resp.StatusCode)
// 	}

// 	var result struct {
// 		ID int `json:"id"`
// 	}

// 	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
// 		return 0, err
// 	}

// 	return result.ID, nil
// }

// // Clean up a test subject after testing
// func cleanupTestSubject(api *registryAPI, name string) error {
// 	// Example implementation to delete a test subject
// 	client := &http.Client{}

// 	url := baseRegistryURL + "/subjects/" + name
// 	req, err := http.NewRequest("DELETE", url, nil)
// 	if err != nil {
// 		return err
// 	}

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("failed to delete subject, status: %d", resp.StatusCode)
// 	}

// 	return nil
// }
