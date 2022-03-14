package gauth

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
)

var testServer *httptest.Server
var testRecorder *httptest.ResponseRecorder
var testContext *gin.Context

func NewClient() *http.Client {
	client := testServer.Client()
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client.Jar = cookieJar
	return client
}
