package httpclient

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// function to stop a request from sending
func stopRequest(req *http.Request) error {
	return errors.New("Request not sent")
}

func TestClient(t *testing.T) {
	assert := assert.New(t)
	c := New()
	assert.NotNil(c.Client)
	assert.NotEmpty(c.Headers["User-Agent"])
}

func TestBeforeRequest(t *testing.T) {
	assert := assert.New(t)
	c := New()
	c.BeforeRequest = stopRequest
	// Request is not actually sent due to stopRequest returning error
	resp, err := c.Get("http://www.example.com")
	assert.Error(err)
	assert.Nil(resp)
}

func TestCopy(t *testing.T) {
	u, _ := url.Parse("http://www.example.com")
	cookies := []*http.Cookie{&http.Cookie{Name: "k", Value: "v"}}
	assert := assert.New(t)
	c := New()
	c.Client.Timeout = 5 * time.Second
	c.Jar.SetCookies(u, cookies)
	c.Headers["h1"] = "v1"
	c.Proxy = "socks5://localhost:9000"

	c2 := c.Copy()
	assert.NotEqual(c, c2)
	assert.NotEqual(c.Client, c2.Client)
	// assert.NotEqual(c.Client.Jar, c2.Client.Jar)

}
