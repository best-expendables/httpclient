package middleware

import (
	"github.com/best-expendables/httpclient/net/profile"
	"net/http"
)

// NetworkProfiler middleware
type NetworkProfiler struct{}

func NewNetworkProfiler() *NetworkProfiler {
	return new(NetworkProfiler)
}

func (NetworkProfiler) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFn(func(request *http.Request) (*http.Response, error) {
		return next.RoundTrip(profile.Observe(request))
	})
}
