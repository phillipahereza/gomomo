package gomomo

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

var (
	mux *http.ServeMux

	ctx = context.TODO()

	client *Client

	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	client = NewClient("", "sandbox", server.URL)
	urlStr, _ := url.Parse(server.URL)
	client.BaseURL = urlStr
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("Request method = %v, expected %v\n", r.Method, expected)
	}
}

type headers map[string]string

func testHeaders(t *testing.T, r *http.Request, expectedHeaders headers) {
	for k, v := range expectedHeaders {
		got := r.Header.Get(k)
		if got != v {
			t.Errorf("Header %s missing", k)
		}
	}
}
