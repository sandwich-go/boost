package httputil

import (
	"io"
	"net/http"
)

// Client HTTP Client 接口
type Client interface {
	// Post 发送 POST 请求
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
	// Get 发送 GET 请求
	Get(url string) (*http.Response, error)
	// Bytes 发送 GET 请求，返回 bytes 数据
	Bytes(url string) ([]byte, error)
	// String 发送 GET 请求，返回 string 数据
	String(url string) (string, error)
	// JSON 发送 GET 请求，返回 JSON 数据
	JSON(url string, v interface{}) error
}

// Error is the custom error type returns from HTTP requests.
type Error struct {
	Message    string
	StatusCode int
	URL        string
}

// Error returns the error message.
func (e *Error) Error() string { return e.Message }
