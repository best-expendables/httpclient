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
		middleware.NewNewrelicApiGateway(middleware.NewURLFormatFunc()),
	)
	return c
}

func NewHttpClientWithMiddlewares(timeout time.Duration, middlewares ...middleware.Middleware) *http.Client {
	return &http.Client{
		Timeout:   timeout,
		Transport: middleware.WithMiddleware(http.DefaultTransport, middlewares...),
	}
}
