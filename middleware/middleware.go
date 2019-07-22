package middleware

import "net/http"

// Middleware for middleware as structure
type Middleware interface {
	RoundTripper(next http.RoundTripper) http.RoundTripper
}

// RoundTripperFn interface for middleware as function
type RoundTripperFn func(request *http.Request) (*http.Response, error)

func (f RoundTripperFn) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

// WithMiddleware general way.
// For flexibility use middleware.Container
func WithMiddleware(rt http.RoundTripper, middlewares ...Middleware) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	for _, middleware := range middlewares {
		rt = middleware.RoundTripper(rt)
	}

	return rt
}
