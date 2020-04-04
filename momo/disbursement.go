package momo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	disbursementsTokenURL           = "/disbursement/token/"
	disbursementsTransferURL        = "/disbursement/v1_0/transfer"
	disbursementsBalanceURL         = "/disbursement/v1_0/account/balance"
	disbursementsIsAccountActiveURl = "/disbursement/v1_0/accountholder/msisdn/"
)

type DisbursementService interface {
	Transfer(ctx context.Context, mobileNumber string, amount int64, id, payeeNote, payerMessage, currency string) (string, error)
	GetTransfer(ctx context.Context, transactionID string) (*PaymentStatusResponse, error)
	GetBalance(ctx context.Context) (*BalanceResponse, error)
	IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error)
	GetToken(ctx context.Context, apiKey, userID string) (string, error)
}

type DisbursementOp struct {
	client *Client
}

func (c *DisbursementOp) GetBalance(ctx context.Context) (*BalanceResponse, error) {
	req, err := c.client.NewRequest(ctx, http.MethodGet, disbursementsBalanceURL, nil)
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

func (c *DisbursementOp) IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error) {
	urlStr := fmt.Sprintf("%s/%s/active", disbursementsIsAccountActiveURl, mobileNumber)
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

func (c *DisbursementOp) GetToken (ctx context.Context, apiKey, userID string) (string, error) {

	req, err := c.client.NewRequest(ctx, http.MethodPost, disbursementsTokenURL, nil)
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(userID, apiKey)

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return "", err
	}
	token := &tokenResponse{}

	err = json.Unmarshal(res.Body, token)
	if err != nil {
		return "", err
	}
	c.client.Token = token.AccessToken
	return token.AccessToken, err
}

func (c *DisbursementOp) Transfer(ctx context.Context, mobileNumber string, amount int64, id, payeeNote, payerMessage, currency string) (string, error) {
	if c.client.Environment == "sandbox" {
		currency = "EUR"
	}

	requestBody := TransferRequestBody{
		Amount:     amount,
		Currency:   currency,
		ExternalID: id,
		Payee: paymentDetails{
			PartyIDType: "MSISDN",
			PartyID:     mobileNumber,
		},
		PayerMessage: payerMessage,
		PayeeNote:    payeeNote,
	}

	req, err := c.client.NewRequest(ctx, http.MethodPost, disbursementsTransferURL, requestBody)
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

func (c *DisbursementOp) GetTransfer(ctx context.Context, transferID string) (*PaymentStatusResponse, error) {
	urlStr := fmt.Sprintf("%s/%s", disbursementsTransferURL, transferID)
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
