package middleware

import (
	"net/http"

	"github.com/best-expendables/trace"
)

// RequestID
type RequestID struct{}

// NewRequestID pass "context-id" into the request
// DEPRECATED should be replaced to Open Trace
func NewRequestID() *RequestID {
	return &RequestID{}
}

func (RequestID) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFn(func(request *http.Request) (*http.Response, error) {
		if requestID := trace.RequestIDFromContext(request.Context()); requestID != "" {
			request = cloneRequest(request)
			trace.RequestIDToHeader(request.Header, requestID)
		}

		return next.RoundTrip(request)
	})
}
