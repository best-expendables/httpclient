package middleware

import (
	"net/http"
)

// Authentication through bearer token
type Authentication struct {
	token string
}

func NewAuthentication(token string) *Authentication {
	return &Authentication{token: token}
}

func (a Authentication) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFn(func(request *http.Request) (*http.Response, error) {
		request = cloneRequest(request)
		request.Header.Set(a.Header())
		return next.RoundTrip(request)
	})
}

// Header returns header name and header value
func (a Authentication) Header() (string, string) {
	return "Authorization", "Bearer " + a.token
}
