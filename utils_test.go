package gauth

import (
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
)

var testServer *httptest.Server

func NewClient() *http.Client {
	client := testServer.Client()
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client.Jar = cookieJar
	return client
}
