package momo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	remittancesTokenURL           = "/remittance/token/"
	remittancesTransferURL        = "/remittance/v1_0/transfer"
	remittancesBalanceURL         = "/remittance/v1_0/account/balance"
	remittancesIsAccountActiveURL = "/remittance/v1_0/accountholder/msisdn/"
)

// RemittanceService handles communication with Remittance related methods of the
// Momo API to remit funds to local recipients from the diaspora
type RemittanceService interface {
	Transfer(ctx context.Context, mobile string, amount int64, id, payeeNote, payerMessage, currency string) (string, error)
	GetTransfer(ctx context.Context, transactionID string) (*PaymentStatusResponse, error)
	GetBalance(ctx context.Context) (*BalanceResponse, error)
	IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error)
	GetToken(ctx context.Context, apiKey, userID string) (string, error)
}

// RemittanceServiceOp handles communication with the Remittance related methods of the Momo API.
type RemittanceServiceOp struct {
	client *Client
}

// GetBalance returns the balance of the account
func (c *RemittanceServiceOp) GetBalance(ctx context.Context) (*BalanceResponse, error) {
	req, err := c.client.NewRequest(ctx, http.MethodGet, remittancesBalanceURL, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code: %d with error %s", res.StatusCode, string(res.Body))
	}

	balance := &BalanceResponse{}
	err = json.Unmarshal(res.Body, balance)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

// IsPayeeActive checks if an account holder is registered and active in the system
func (c *RemittanceServiceOp) IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error) {
	urlStr := fmt.Sprintf("%s/%s/active", remittancesIsAccountActiveURL, mobileNumber)
	req, err := c.client.NewRequest(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return false, err
	}

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return false, err
	}

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("response code: %d with error %s", res.StatusCode, string(res.Body))
	}

	return true, nil
}

// GetToken creates an access token which can then be used to authorize and authenticate towards the other end-points of the Remittance API
func (c *RemittanceServiceOp) GetToken(ctx context.Context, apiKey, userID string) (string, error) {

	req, err := c.client.NewRequest(ctx, http.MethodPost, remittancesTokenURL, nil)
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
	return token.AccessToken, nil
}

// Transfer operation is used to transfer an amount from the owner’s account to a payee account.
func (c *RemittanceServiceOp) Transfer(ctx context.Context, mobile string, amount int64, id, payeeNote, payerMessage, currency string) (string, error) {
	if c.client.Environment == "sandbox" {
		currency = "EUR"
	}

	requestBody := transferRequestBody{
		Amount:     amount,
		Currency:   currency,
		ExternalID: id,
		Payee: paymentDetails{
			PartyIDType: "MSISDN",
			PartyID:     mobile,
		},
		PayerMessage: payerMessage,
		PayeeNote:    payeeNote,
	}

	req, err := c.client.NewRequest(ctx, http.MethodPost, remittancesTransferURL, requestBody)
	if err != nil {
		return "", err
	}

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("response code: %d with error %s", res.StatusCode, string(res.Body))
	}

	return req.Header.Get("X-Reference-Id"), nil
}

// GetTransfer retrieves transfer information using the transactionId returned by Transfer
func (c *RemittanceServiceOp) GetTransfer(ctx context.Context, transferID string) (*PaymentStatusResponse, error) {
	urlStr := fmt.Sprintf("%s/%s", remittancesTransferURL, transferID)
	req, err := c.client.NewRequest(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(ctx, req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response code: %d with error %s", res.StatusCode, string(res.Body))
	}

	status := &PaymentStatusResponse{}
	err = json.Unmarshal(res.Body, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}
