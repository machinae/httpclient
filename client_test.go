package httpclient

import (
	"errors"
	"net/http"
	"testing"

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
