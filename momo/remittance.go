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
	RemittancesTokenURL           = "/remittance/token/"
	RemittancesTransferURL        = "/remittance/v1_0/transfer"
	RemittancesBalanceURL         = "/remittance/v1_0/account/balance"
	RemittancesIsAccountActiveURl = "/remittance/v1_0/accountholder/msisdn/"
)

type RemittanceService interface {
	Transfer(ctx context.Context, mobile string, amount int64, id, payeeNote, payerMessage, currency string) (string, error)
	GetTransfer(ctx context.Context, transactionID string) (*PaymentStatusResponse, error)
	GetBalance(ctx context.Context) (*BalanceResponse, error)
	IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error)
	GetToken(ctx context.Context) error
}

type RemittanceOp struct {
	client *momoClient
}

func (c *RemittanceOp) GetBalance(ctx context.Context) (*BalanceResponse, error) {
	req, err := c.client.NewRequest(ctx, http.MethodGet, RemittancesBalanceURL, nil)
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

func (c *RemittanceOp) IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error) {
	urlStr := fmt.Sprintf("%s/%s/active", RemittancesIsAccountActiveURl, mobileNumber)
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

func (c *RemittanceOp) GetToken(ctx context.Context) error {
	apiKey := os.Getenv("COLLECTION_API_KEY")
	userID := os.Getenv("COLLECTION_USER_ID")
	if userID == "" {
		return errors.New("COLLECTION_USER_ID should be set in the environment")
	}

	if apiKey == "" {
		return errors.New("COLLECTION_API_KEY should be set in the environment")
	}

	req, err := c.client.NewRequest(ctx, http.MethodPost, RemittancesTokenURL, nil)
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

func (c *RemittanceOp) Transfer(ctx context.Context, mobile string, amount int64, id, payeeNote, payerMessage, currency string) (string, error) {
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

	req, err := c.client.NewRequest(ctx, http.MethodPost, RemittancesTransferURL, requestBody)
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

func (c *RemittanceOp) GetTransfer(ctx context.Context, transferID string) (*PaymentStatusResponse, error) {
	urlStr := fmt.Sprintf("%s/%s", RemittancesTransferURL, transferID)
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

func NewRemittanceClient(key, environment, baseURL string) *RemittanceOp {
	c := newClient(key, environment, baseURL)
	return &RemittanceOp{client: c}
}
