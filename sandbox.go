package gomomo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// SandboxService handles communication with sandbox related methods of the Momo API
type SandboxService interface {
	CreateSandboxUser(callbackHost string) (string, error)
	GenerateSandboxUserAPIKey(referenceID string) (*APIKeyResponse, error)
}

// SandboxServiceOp handles communication with methods on Momo API to create Sandbox users
type SandboxServiceOp struct {
	client *Client
}

// APIKeyResponse structure for returning API Key
type APIKeyResponse struct {
	APIKey string `json:"apiKey"`
}

// CreateSandboxUser creates a user to test the Momo APU in a sandbox environment
func (c *SandboxServiceOp) CreateSandboxUser(callbackHost string) (string, error) {
	ctx := context.Background()
	body := map[string]string{
		"providerCallbackHost": callbackHost,
	}
	req, err := c.client.NewRequest(ctx, http.MethodPost, "v1_0/apiuser", body)
	if err != nil {
		return "", err
	}
	response, err := c.client.Do(ctx, req)
	if err != nil {
		return "", err
	}

	if response.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("response code: %d with error %s", response.StatusCode, string(response.Body))
	}
	return response.ReferenceID, nil
}

// GenerateSandboxUserAPIKey is used to create an API key for an API user in the sandbox target environment
func (c *SandboxServiceOp) GenerateSandboxUserAPIKey(referenceID string) (*APIKeyResponse, error) {
	urlStr := fmt.Sprintf("v1_0/apiuser/%s/apikey", referenceID)
	ctx := context.Background()
	req, err := c.client.NewRequest(ctx, http.MethodPost, urlStr, nil)
	if err != nil {
		return nil, err
	}
	response, err := c.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("response code: %d with error %s", response.StatusCode, string(response.Body))
	}

	keyResponse := &APIKeyResponse{}
	err = json.Unmarshal(response.Body, keyResponse)
	if err != nil {
		return nil, err
	}
	return keyResponse, nil
}
