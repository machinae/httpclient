package httpclient

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	assert := assert.New(t)
	c := New()
	assert.NotNil(c.Client)
	assert.NotEmpty(c.Headers["User-Agent"])
}
