package momo

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
	mediaType  = "application/json"
)

type Client struct {
	client          *http.Client
	BaseURL         *url.URL
	SubscriptionKey string
	Token           string
	Environment     string
}

type Response struct {
	StatusCode int
	Body       []byte
	Headers    map[string][]string
	ReferenceID string
}


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

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("X-Reference-Id", uuid.New().String())
	req.Header.Add("Ocp-Apim-Subscription-Key", c.SubscriptionKey)
	if c.Token != "" {
		req.Header.Add("Authorization", "Bearer " + c.Token)
	}

	return req, nil
}

func (c *Client) Do(ctx context.Context, req *http.Request) (*Response, error) {
	req.WithContext(ctx)
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	response, err := BuildResponse(res)
	if err != nil {
		return nil, err
	}
	response.ReferenceID = req.Header.Get("X-Reference-Id")
	return response, nil
}

func BuildResponse(res *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(res.Body)
	response := Response{
		StatusCode: res.StatusCode,
		Body:       body,
		Headers:    res.Header,
		ReferenceID: "",
	}
	res.Body.Close()
	return &response, err
}

func NewClient(key, environment, baseUrl string) *Client {
	urlStr, err := url.Parse(baseUrl)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		client:          http.DefaultClient,
		BaseURL:         urlStr,
		SubscriptionKey: key,
		Token:           "",
		Environment:     environment,
	}
}