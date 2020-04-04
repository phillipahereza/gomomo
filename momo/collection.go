package momo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	collectionsTokenURL           = "/collection/token/"
	collectionsRequestToPayURL    = "/collection/v1_0/requesttopay"
	collectionsBalanceURL         = "/collection/v1_0/account/balance"
	collectionsIsAccountActiveURl = "/collection/v1_0/accountholder/msisdn/"
)

type CollectionService interface {
	RequestToPay(ctx context.Context, mobile string, amount int64, id, payeeNote, payerMessage, currency string) (string, error)
	GetTransaction(ctx context.Context, transactionID string) (*PaymentStatusResponse, error)
	GetBalance(ctx context.Context) (*BalanceResponse, error)
	IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error)
	GetToken(ctx context.Context, apiKey, userID string) (string, error)
}

type CollectionServiceOp struct {
	client *Client
}

var _ CollectionService= &CollectionServiceOp{}

func (c *CollectionServiceOp) RequestToPay(ctx context.Context, mobile string, amount int64, id, payeeNote, payerMessage, currency string) (string, error) {
	if c.client.Environment == "sandbox" {
		currency = "EUR"
	}

	requestBody := PaymentRequestBody{
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
	transactionID := req.Header.Get("X-Reference-Id")

	if err != nil {
		return "", err
	}

	res, err := c.client.Do(ctx, req)
	if err != nil {
		log.Println(string(res.Body))
		return "", err
	}

	if res.StatusCode != http.StatusAccepted {
		return "", errors.New(fmt.Sprintf("response code: %d with error %s", res.StatusCode, string(res.Body)))
	}

	return transactionID, nil
}

func (c *CollectionServiceOp) GetTransaction(ctx context.Context, transactionID string) (*PaymentStatusResponse, error) {
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

	status := &PaymentStatusResponse{}
	err = json.Unmarshal(res.Body, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (c *CollectionServiceOp) GetBalance(ctx context.Context) (*BalanceResponse, error) {
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

func (c *CollectionServiceOp) IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error) {
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

func (c *CollectionServiceOp) GetToken(ctx context.Context, apiKey, userID string) (string, error) {
	req, err := c.client.NewRequest(ctx, http.MethodPost, collectionsTokenURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(userID, apiKey)

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		code := strconv.Itoa(res.StatusCode)
		return "", errors.New("status Code received is " + code)
	}
	token := &tokenResponse{}

	err = json.Unmarshal(res.Body, token)
	if err != nil {
		return "", err
	}
	c.client.Token = token.AccessToken
	return token.AccessToken, nil
}
