package momo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type CollectionsOp struct {
	client *momoClient
}

func (c *CollectionsOp) RequestToPay(mobile string, amount int64, id, payeeNote, payerMessage, currency string) (string) {
	return "transaction ID"
}

func (c *CollectionsOp) GetTransaction(transactionID string) {
}

func (c *CollectionsOp) GetBalance() {}

func (c *CollectionsOp) IsPayeeActive(accountHolderType, accountHolderID string) {

}

func (c *CollectionsOp) GetToken() (error) {
	apiKey := os.Getenv("COLLECTION_API_KEY")
	userID := os.Getenv("COLLECTION_USER_ID")
	if userID == "" {
		return errors.New("COLLECTION_USER_ID should be set in the environment")
	}

	if apiKey == "" {
		return errors.New("COLLECTION_API_KEY should be set in the environment")
	}

	ctx := context.Background()

	req, err := c.client.NewRequest(ctx, http.MethodPost, "collection/token/", nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(userID, apiKey)

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return err
	}
	token := &tokenResponse{}

	err = json.Unmarshal(res.Body, token)
	if err != nil {
		return err
	}
	c.client.Token = token.AccessToken
	return nil
}

func NewCollectionsClient(key, environment, baseURL string) *CollectionsOp {
	c := newClient(key, environment, baseURL)
	return &CollectionsOp{client: c}
}
