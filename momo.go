package gomomo

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	mediaType = "application/json"
)

// Client manages communication with MTN Momo API.
type Client struct {
	client          *http.Client
	BaseURL         *url.URL
	SubscriptionKey string
	Token           string
	Environment     string
	Collection      CollectionService
	Disbursement    DisbursementService
	Remittance      RemittanceService
	Sandbox         SandboxService
}

// Response returned by API calls
type Response struct {
	StatusCode  int
	Body        []byte
	Headers     map[string][]string
	ReferenceID string
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

// BalanceResponse holds the account Balance
type BalanceResponse struct {
	AvailableBalance string `json:"availableBalance"`
	Currency         string `json:"currency"`
}

type paymentDetails struct {
	PartyIDType string `json:"partyIdType"`
	PartyID     string `json:"partyId"`
}

type paymentRequestBody struct {
	Amount       int64          `json:"amount"`
	Currency     string         `json:"currency"`
	ExternalID   string         `json:"externalId"`
	Payer        paymentDetails `json:"payer"`
	PayerMessage string         `json:"payerMessage"`
	PayeeNote    string         `json:"payeeNote"`
}

type transferRequestBody struct {
	Amount       int64          `json:"amount"`
	Currency     string         `json:"currency"`
	ExternalID   string         `json:"externalId"`
	Payee        paymentDetails `json:"payee"`
	PayerMessage string         `json:"payerMessage"`
	PayeeNote    string         `json:"payeeNote"`
}

type reason struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// PaymentStatusResponse returned for every successful call to make a transfer
type PaymentStatusResponse struct {
	Amount                 string         `json:"amount,omitempty"`
	Currency               string         `json:"currency,omitempty"`
	FinancialTransactionID int64          `json:"financialTransactionId,omitempty"`
	ExternalID             string         `json:"externalId,omitempty"`
	Payer                  paymentDetails `json:"paymentDetails,omitempty"`
	Status                 string         `json:"status,omitempty"`
	Reason                 string         `json:"reason,omitempty"`
}

// NewRequest creates an API request. A relative URL can be provided in urlStr, which will be resolved to the
// BaseURL of the Client.
func (c *Client) NewRequest(ctx context.Context, method, urlStr string, body interface{}) (*http.Request, error) {
	u, err := c.BaseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if body != nil {
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.WithContext(ctx)

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("X-Reference-Id", uuid.New().String())
	req.Header.Add("Ocp-Apim-Subscription-Key", c.SubscriptionKey)

	if c.Environment != "" {
		req.Header.Add("X-Target-Environment", c.Environment)
	}
	if c.Token != "" {
		req.Header.Add("Authorization", "Bearer "+c.Token)
	}

	return req, nil
}

// Do sends an API request and returns the API response.
func (c *Client) Do(ctx context.Context, req *http.Request) (*Response, error) {
	req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	response, err := buildResponse(res)
	if err != nil {
		return nil, err
	}
	response.ReferenceID = req.Header.Get("X-Reference-Id")
	return response, nil
}

func buildResponse(res *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	response := Response{
		StatusCode:  res.StatusCode,
		Body:        body,
		Headers:     res.Header,
		ReferenceID: "",
	}
	res.Body.Close()
	return &response, err
}

// NewClient returns a new Momo API client, using the given
// http.Client to perform all requests.
func NewClient(key, environment, baseURL string) *Client {
	urlStr, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	c := &Client{
		client:          http.DefaultClient,
		BaseURL:         urlStr,
		SubscriptionKey: key,
		Token:           "",
		Environment:     environment,
	}
	c.Collection = &CollectionServiceOp{client: c}
	c.Disbursement = &DisbursementServiceOp{client: c}
	c.Remittance = &RemittanceServiceOp{client: c}
	return c
}
