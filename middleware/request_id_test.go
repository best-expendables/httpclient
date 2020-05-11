package middleware

import (
	"github.com/best-expendables/trace"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if requestID := trace.RequestIDFromHeader(r.Header); requestID == "" {
			t.Fatal("Empty context-id")
		} else if requestID != "TEST-REQUEST-ID" {
			t.Errorf("Request ID not equals. Expected 'TEST-REQUEST-ID', actual %s", requestID)
		}

	}))
	defer srv.Close()

	client := &http.Client{
		Transport: WithMiddleware(nil, NewRequestID()),
	}

	request, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	ctx := trace.ContextWithRequestID(request.Context(), "TEST-REQUEST-ID")
	client.Do(request.WithContext(ctx))
}
