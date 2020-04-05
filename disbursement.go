package gomomo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	disbursementsTokenURL           = "/disbursement/token/"
	disbursementsTransferURL        = "/disbursement/v1_0/transfer"
	disbursementsBalanceURL         = "/disbursement/v1_0/account/balance"
	disbursementsIsAccountActiveURL = "/disbursement/v1_0/accountholder/msisdn/"
)

// DisbursementService handles communication with Disbursement related methods of the
// Momo API to automatically deposit funds into multiple users accounts
type DisbursementService interface {
	Transfer(ctx context.Context, mobileNumber string, amount int64, id, payeeNote, payerMessage, currency string) (string, error)
	GetTransfer(ctx context.Context, transactionID string) (*PaymentStatusResponse, error)
	GetBalance(ctx context.Context) (*BalanceResponse, error)
	IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error)
	GetToken(ctx context.Context, apiKey, userID string) (string, error)
}

// DisbursementServiceOp handles communication with the Disbursement related methods of the Momo API.
type DisbursementServiceOp struct {
	client *Client
}

// GetBalance returns the balance of the account
func (c *DisbursementServiceOp) GetBalance(ctx context.Context) (*BalanceResponse, error) {
	req, err := c.client.NewRequest(ctx, http.MethodGet, disbursementsBalanceURL, nil)
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
func (c *DisbursementServiceOp) IsPayeeActive(ctx context.Context, mobileNumber string) (bool, error) {
	urlStr := fmt.Sprintf("%s/%s/active", disbursementsIsAccountActiveURL, mobileNumber)
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

// GetToken creates an access token which can then be used to authorize and authenticate towards the other end-points of the Disbursement API
func (c *DisbursementServiceOp) GetToken(ctx context.Context, apiKey, userID string) (string, error) {

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

// Transfer operation is used to transfer an amount from the owner’s account to a payee account.
func (c *DisbursementServiceOp) Transfer(ctx context.Context, mobileNumber string, amount int64, id, payeeNote, payerMessage, currency string) (string, error) {
	if c.client.Environment == "sandbox" {
		currency = "EUR"
	}

	requestBody := transferRequestBody{
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
		return "", fmt.Errorf("response code: %d with error %s", res.StatusCode, string(res.Body))
	}

	return req.Header.Get("X-Reference-Id"), nil
}

// GetTransfer retrieves transfer information using the transactionId returned by Transfer
func (c *DisbursementServiceOp) GetTransfer(ctx context.Context, transferID string) (*PaymentStatusResponse, error) {
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
		return nil, fmt.Errorf("response code: %d with error %s", res.StatusCode, string(res.Body))
	}

	status := &PaymentStatusResponse{}
	err = json.Unmarshal(res.Body, status)
	if err != nil {
		return nil, err
	}
	return status, nil
}
