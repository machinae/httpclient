// Package httpclient provides a helpful wrapper around the standard http.Client
// to make it behave more like a web browser
package httpclient

import (
	"errors"
	"net/http"
	"net/url"
	"time"
)

var (
	// Default timeout for requests
	Timeout = 60 * time.Second

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
		Proxy:               c.proxy,
		MaxIdleConnsPerHost: MaxConnsPerHost,
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

	return c.Client.Do(req)
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
