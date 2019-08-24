package middleware

import (
	"bitbucket.org/snapmartinc/httpclient/net/profile"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNetworkProfiler(t *testing.T) {
	a := assert.New(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Nanosecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := &http.Client{
		Transport: WithMiddleware(nil, NewNetworkProfiler()),
	}

	req, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	res, err := client.Do(req)
	a.NoError(err)

	report := profile.ReportFromResponse(res)
	a.NotNil(report)
	a.NotZero(report.ConnectionTime())
}
