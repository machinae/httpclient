// Package httpclient provides a helpful wrapper around the standard http.Client
// to make it behave more like a web browser
package httpclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// Default timeout for requests
	Timeout = 30 * time.Second

	// Default headers for all requests
	Headers = map[string]string{
		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8",
		"Accept-Language":           "en-US,en;q=0.9",
		"Upgrade-Insecure-Requests": "1",
		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36",
	}

	// Maximum number of persistent connections per server
	MaxConnsPerHost = 10
)

// An HTTP client that behaves more like a web browser
type Client struct {
	// The underlying Client
	*http.Client
	// Persistent headers to send with every request
	Headers map[string]string

	// URL of a proxy to use for all requests
	Proxy string

	// If not nil, BeforeRequest is called on each Request before it is
	// sent. If this function returns an error, the request is not sent.
	// This can be used to globally transform or log all requests before sending
	// them, or to filter which requests get sent
	BeforeRequest func(*http.Request) error
}

// Initializes a new http client
func New() *Client {
	c := &Client{
		Client: &http.Client{
			Timeout: Timeout,
			Jar:     newCookieJar(),
		},
		Headers: make(map[string]string),
	}
	for k, v := range Headers {
		c.Headers[k] = v
	}

	tr := &http.Transport{
		Proxy: c.proxy,
		DialContext: (&net.Dialer{
			Timeout:   Timeout,
			KeepAlive: Timeout,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   MaxConnsPerHost,
		IdleConnTimeout:       3 * Timeout,
		TLSHandshakeTimeout:   Timeout,
		ExpectContinueTimeout: 1 * time.Second,
	}
	c.Transport = tr
	return c
}

// Returns the current proxy
func (c *Client) proxy(req *http.Request) (*url.URL, error) {
	if c.Proxy == "" {
		return nil, nil
	}
	return url.Parse(c.Proxy)
}

// Overrides of Client methods

// Do sends an HTTP request and returns an HTTP response
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	req, err := c.prepareRequest(req)
	if err != nil {
		return nil, err
	}

	if c.BeforeRequest != nil {
		err := c.BeforeRequest(req)
		if err != nil {
			return nil, err
		}
	}

	return c.Client.Do(req)
}

// Make a GET request
func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Make a HEAD request
func (c *Client) Head(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Make a POST request
func (c *Client) Post(url string, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

// Make a POST request with encoded form values
func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

// Make a POST request with encoded JSON body
func (c *Client) PostJson(url string, data interface{}) (resp *http.Response, err error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return c.Post(url, "application/json", bytes.NewReader(body))
}

// Wraps request with additions from the client
func (c *Client) prepareRequest(req *http.Request) (*http.Request, error) {
	if req == nil {
		return nil, errors.New("Request is empty")
	}
	c.setHeaders(req)
	return req, nil
}

// Merges request headers with those defined on client
// Existing request headers are not overwritten
func (c *Client) setHeaders(req *http.Request) {
	for k, v := range c.Headers {
		k = http.CanonicalHeaderKey(k)
		if rh := req.Header.Get(k); len(rh) == 0 {
			req.Header.Set(k, v)
		}
	}
}
