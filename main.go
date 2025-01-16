package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Terraform Cloud API base URL
const baseURL = "https://app.terraform.io/api/v2"

// Configuration struct for Terraform Cloud
type Config struct {
	Token       string // Terraform Cloud API Token
	PolicySetID string // Policy Set ID
	PolicyDir   string // Directory containing .rego files
}

func GetPolicyId(config Config, policyName string) (string, error) {
	// Construct the URL
	url := fmt.Sprintf("%s/organizations/%s/policies", baseURL, "DJB-Personal")

	// Create a new request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/vnd.api+json")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s", string(body))
	}

	// Print the data
	var result map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	// After your existing JSON unmarshal code
	if data, ok := result["data"].([]interface{}); ok {
		for _, item := range data {
			// Cast item to map to access its properties
			if policy, ok := item.(map[string]interface{}); ok {
				// Now you can access fields within each policy
				id := policy["id"]
				attributes := policy["attributes"].(map[string]interface{})
				name := attributes["name"]

				if name == policyName {
					return id.(string), nil
				}
			}
		}
	} else {
		fmt.Println("Could not parse data array")
	}

	return "Policy DNE", nil
}

func UploadPolicy(config Config, policyContent []byte, policyId string) error {
	// Construct the URL
	url := fmt.Sprintf("%s/policies/%s/upload", baseURL, policyId)

	// Create a new request
	req, err := http.NewRequest("PUT", url, io.NopCloser(bytes.NewReader(policyContent)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "Bearer "+config.Token)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s", string(body))
	}

	return nil
}

// Main function
func main() {
	token := os.Getenv("TF_API_TOKEN")
	if token == "" {
		log.Fatal("TF_API_TOKEN environment variable is not set")
	}

	policySetId := os.Getenv("POLICY_SET_ID")
	if policySetId == "" {
		log.Fatal("POLICY_SET_ID environment variable is not set")
	}

	config := Config{
		Token:       token,
		PolicySetID: policySetId,
		PolicyDir:   "./policies", // Directory containing your .rego files
	}

	// Iterate through all .rego files in the policy directory
	files, err := ioutil.ReadDir(config.PolicyDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading policy directory: %v\n", err)
		os.Exit(1)
	}

	var policyIDs []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".rego" {
			policyContent, err := ioutil.ReadFile(filepath.Join(config.PolicyDir, file.Name()))

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error reading policy file %s: %v\n", file.Name(), err)
				continue
			}

			// Get the Policy ID
			policyID, err := GetPolicyId(config, file.Name())
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting policy ID for %s: %v\n", file.Name(), err)
				continue
			}

			err = UploadPolicy(config, policyContent, policyID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error uploading policy %s: %v\n", file.Name(), err)
				continue
			}

			policyIDs = append(policyIDs, policyID)
		}
	}

	fmt.Printf("The following policies have been updated: %v", policyIDs)
}
