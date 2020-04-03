package momo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

const (
	collectionsRequestToPayURL = "/collection/v1_0/requesttopay"
	collectionsBalanceURL      = "/collection/v1_0/account/balance"
	collectionsIsAccountActiveURl = "/collection/v1_0/accountholder/msisdn/"
)

type CollectionsOp struct {
	client *momoClient
}

type payer struct {
	PartyIDType string `json:"partyIdType"`
	PartyID     string `json:"partyId"`
}

type requestToPayBody struct {
	Amount       int64  `json:"amount"`
	Currency     string `json:"currency"`
	ExternalID   string `json:"externalId"`
	Payee        payer  `json:"payer"`
	PayerMessage string `json:"payerMessage"`
	PayeeNote    string `json:"payeeNote"`
}

type Reason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type CollectionStatusResponse struct {
	Amount                 string `json:"amount,omitempty"`
	Currency               string `json:"amount,omitempty"`
	FinancialTransactionID string `json:"financialTransactionId,omitempty"`
	ExternalID             string `json:"externalId,omitempty"`
	Payer                  payer  `json:"payer,omitempty"`
	Status                 string `json:"status,omitempty"`
	Reason                 Reason `json:"reason,omitempty"`
}

func (c *CollectionsOp) RequestToPay(ctx context.Context, mobile string, amount int64, id, payeeNote, payerMessage, currency string) (string, error) {
	if c.client.Environment == "sandbox" {
		currency = "EUR"
	}

	requestBody := requestToPayBody{
		Amount:     amount,
		Currency:   currency,
		ExternalID: id,
		Payee: payer{
			PartyIDType: "MSISDN",
			PartyID:     mobile,
		},
		PayerMessage: payerMessage,
		PayeeNote:    payeeNote,
	}

	req, err := c.client.NewRequest(ctx, http.MethodPost, collectionsRequestToPayURL, requestBody)
	if err != nil {
		return "", err
	}

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusAccepted {
		return "", errors.New(fmt.Sprintf("response code: %d with error %s", res.StatusCode, string(res.Body)))
	}

	return req.Header.Get("X-Reference-Id"), nil
}

func (c *CollectionsOp) GetTransaction(ctx context.Context, transactionID string) (*CollectionStatusResponse, error) {
	urlStr := fmt.Sprintf("%s/%s", collectionsRequestToPayURL, transactionID)
	req, err := c.client.NewRequest(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response code: %d with error %s", res.StatusCode, string(res.Body)))
	}

	status := &CollectionStatusResponse{}
	err = json.Unmarshal(res.Body, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (c *CollectionsOp) GetBalance(ctx context.Context) (*BalanceResponse, error) {
	req, err := c.client.NewRequest(ctx, http.MethodGet, collectionsBalanceURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("response code: %d with error %s", res.StatusCode, string(res.Body)))
	}

	balance := &BalanceResponse{}
	err = json.Unmarshal(res.Body, balance)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (c *CollectionsOp) IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error) {
	urlStr := fmt.Sprintf("%s/%s/active", collectionsIsAccountActiveURl, mobileNumber)
	req, err := c.client.NewRequest(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return false, err
	}

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return false, err
	}

	if res.StatusCode != http.StatusOK {
		return false, errors.New(fmt.Sprintf("response code: %d with error %s", res.StatusCode, string(res.Body)))
	}

	return true, nil
}

func (c *CollectionsOp) GetToken(ctx context.Context) error {
	apiKey := os.Getenv("COLLECTION_API_KEY")
	userID := os.Getenv("COLLECTION_USER_ID")
	if userID == "" {
		return errors.New("COLLECTION_USER_ID should be set in the environment")
	}

	if apiKey == "" {
		return errors.New("COLLECTION_API_KEY should be set in the environment")
	}

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
