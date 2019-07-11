package mtnmomo

import "net/http"

type MomoClient struct {
	Environment  string
	BaseURL      string
	CallbackHost string

	CollectionPrimaryKey string
	CollectionUserID     string
	CollectionAPISecret  string

	RemittancePrimaryKey string
	RemittanceUserID     string
	RemittanceAPISecret  string

	DisbursementPrimaryKey string
	DisbursementUserID     string
	DisbursementAPISecret  string

	client http.Client
}

type response interface {
}

type MomoProduct interface {
	getAuthToken(url, subscriptionKey string) string
	getBalance(url, subscriptionKey string) response
	getTransactionStatus(transactionID, url, subscriptionKey string) string
}
