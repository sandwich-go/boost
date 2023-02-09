package httputil

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// A Client is an HTTP client.
// It wraps net/http's client and add some methods for making HTTP request easier.
type httpClient struct {
	*http.Client
}

// New create new an HTTP client.
func New(c *http.Client) Client { return &httpClient{Client: c} }

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
	return ioutil.ReadAll(resp.Body)
}

func (c *httpClient) String(url string) (string, error) {
	bytes, err := c.Bytes(url)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

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
