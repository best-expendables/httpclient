package httpclient

import (
	"context"
	"testing"
	"time"

	"github.com/best-expendables/logger"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	c := NewDefaultHttpClient(logger.EntryFromContext(context.Background()), 2*time.Second)

	asserts := assert.New(t)
	asserts.NotNil(c)
	asserts.Equal(2*time.Second, c.Timeout)
}
