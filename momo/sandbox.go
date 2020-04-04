package momo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type SandboxOp struct {
	client *Client
}

type SandboxService interface {
	CreateSandboxUser(callbackHost string) (string, error)
	GenerateSandboxUserAPIKey(referenceId string) (*APIKeyResponse, error)
}

type APIKeyResponse struct {
	APIKey string `json:"apiKey"`
}

// create  user
func (c *SandboxOp) CreateSandboxUser(callbackHost string) (string, error) {
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
		return "", errors.New(fmt.Sprintf("returned with status %d", response.StatusCode))
	}
	return response.ReferenceID, nil
}

// create user API key
func (c *SandboxOp) GenerateSandboxUserAPIKey(referenceId string) (*APIKeyResponse, error) {
	urlStr := fmt.Sprintf("v1_0/apiuser/%s/apikey", referenceId)
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
		return nil, errors.New(fmt.Sprintf("returned with status %d", response.StatusCode))
	}

	keyResponse := &APIKeyResponse{}
	err = json.Unmarshal(response.Body, keyResponse)
	if err != nil {
		return nil, err
	}
	return keyResponse, nil
}

