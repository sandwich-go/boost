package httputil

import (
	"crypto/tls"
	"net/http"
	"time"
)

var defaultHTTPTimeout = time.Duration(10) * time.Second
var defaultDNSCacheTime = time.Minute

// A Client is an HTTP client.
// It wraps net/http's client and add some methods for making HTTP request easier.
type httpClient struct {
	*http.Client
}

var globalClient httpClient

func SetHTTPClient(c *http.Client) { globalClient.Client = c }

func init() {
	dc := NewDNSCache(defaultDNSCacheTime)
	client := &http.Client{
		Timeout: defaultHTTPTimeout,
		Transport: &http.Transport{
			DialContext:            dc.GetDialContext(),
			MaxIdleConns:           50,
			IdleConnTimeout:        60 * time.Second,
			TLSHandshakeTimeout:    5 * time.Second,
			ExpectContinueTimeout:  1 * time.Second,
			MaxResponseHeaderBytes: 5 * 1024,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 单向验证
			},
		},
	}
	SetHTTPClient(client)
}
