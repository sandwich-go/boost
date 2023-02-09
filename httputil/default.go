package httputil

import (
	"crypto/tls"
	"github.com/sandwich-go/boost/httputil/dns"
	"io"
	"net/http"
	"time"
)

var (
	globalClient        httpClient
	defaultHTTPTimeout  = time.Duration(10) * time.Second
	defaultDNSCacheTime = time.Minute
)

// SetDefaultHTTPClient 设置默认的 HTTP Client
func SetDefaultHTTPClient(c *http.Client) { globalClient.Client = c }

// SetDefaultTimeout 设置默认的超时时间
func SetDefaultTimeout(timeout time.Duration) { globalClient.Client.Timeout = timeout }

// Post 使用默认的 Client 发送 POST 请求
func Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return globalClient.Post(url, contentType, body)
}

// Get 使用默认的 Client 发送 GET 请求
func Get(url string) (*http.Response, error) {
	return globalClient.Get(url)
}

// Bytes 使用默认的 Client 发送 GET 请求，返回 bytes 数据
func Bytes(url string) ([]byte, error) {
	return globalClient.Bytes(url)
}

// String 使用默认的 Client 发送 GET 请求，返回 string 数据
func String(url string) (string, error) {
	return globalClient.String(url)
}

// JSON 使用默认的 Client 发送 GET 请求，返回 JSON 数据
func JSON(url string, v interface{}) error {
	return globalClient.JSON(url, v)
}

func init() {
	dc := dns.NewCache(dns.WithTTL(defaultDNSCacheTime))
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
	SetDefaultHTTPClient(client)
}
