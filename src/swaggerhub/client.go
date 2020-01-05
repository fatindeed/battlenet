package swaggerhub

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
)

// NewRequest function
// https://app.swaggerhub.com/apis-docs/swagger-hub/registry-api/1.0.47
func NewRequest(path string, body []byte) (map[string]interface{}, error) {
	req, err := http.NewRequest(
		"POST",
		"https://api.swaggerhub.com"+path,
		bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	// Set headers
	req.Header.Set("Authorization", os.Getenv("SWAGGERHUB_API_KEY"))
	req.Header.Set("Content-Type", "application/yaml")
	req.Header.Set("Accept", "application/json")
	// Set client timeout
	client := &http.Client{}
	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// fmt.Println(resp.Body)
	if len, err := strconv.Atoi(resp.Header["Content-Length"][0]); err == nil && len > 0 {
		var result map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()
		return result, nil
	}
	return nil, nil
}

// SaveDefinition function
func SaveDefinition(owner string, api string, body []byte) error {
	resp, err := NewRequest("/apis/"+owner+"/"+api, body)
	if err != nil {
		return err
	}
	if code, ok := resp["code"]; ok && code.(float64) > 0 {
		return fmt.Errorf("swaggerhub response error: %s", resp["message"])
	}
	return nil
}
