package httpclient

import (
	"net/http"
	"time"

	"bitbucket.org/snapmartinc/httpclient/middleware"
	"bitbucket.org/snapmartinc/logger"
)

func NewDefaultHttpClient(defaultEntry logger.Entry, timeout time.Duration) *http.Client {
	c := &http.Client{
		Timeout: timeout,
	}

	c.Transport = middleware.WithMiddleware(
		c.Transport,
		middleware.NewResponseLogger(defaultEntry),
		middleware.NewRequestLogger(defaultEntry),
		middleware.NewNewrelicApiGateway(),
	)

	return c
}
