package httputil

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

// Error is the custom error type returns from HTTP requests.
type Error struct {
	Message    string
	StatusCode int
	URL        string
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Message
}

func (c *httpClient) err(resp *http.Response, message string) error {
	if message == "" {
		message = fmt.Sprintf("Get %s -> %d", resp.Request.URL.String(), resp.StatusCode)
	}
	return &Error{
		Message:    message,
		StatusCode: resp.StatusCode,
		URL:        resp.Request.URL.String(),
	}
}

// Bytes fetches the specified url and returns the response body as bytes.
func (c *httpClient) Bytes(url string) ([]byte, error) {
	resp, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, c.err(resp, "")
	}
	p, err := ioutil.ReadAll(resp.Body)
	return p, err
}

// String fetches the specified URL and returns the response body as a string.
func (c *httpClient) String(url string) (string, error) {
	bytes, err := c.Bytes(url)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// JSON issues a GET request to a specified URL and unmarshal json data from the response body.
func (c *httpClient) JSON(url string, v interface{}) error {
	resp, err := c.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		return c.err(resp, "")
	}
	err = json.NewDecoder(resp.Body).Decode(v)
	if _, ok := err.(*json.SyntaxError); ok {
		err = c.err(resp, "JSON syntax error at "+url)
	}
	return err
}

// SetDefaultTimeout set timeout of global client
func SetDefaultTimeout(timeout time.Duration) {
	globalClient.Client.Timeout = timeout
}

func Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	return globalClient.Post(url, contentType, body)
}

// Get issues a GET to the specified URL. It returns an http.response for further processing.
func Get(url string) (*http.Response, error) {
	return globalClient.Get(url)
}

// Bytes fetches the specified url and returns the response body as bytes.
func Bytes(url string) ([]byte, error) {
	return globalClient.Bytes(url)
}

// String fetches the specified URL and returns the response body as a string.
func String(url string) (string, error) {
	return globalClient.String(url)
}

// JSON issues a GET request to a specified URL and unmarshal json data from the response body.
func JSON(url string, v interface{}) error {
	return globalClient.JSON(url, v)
}
