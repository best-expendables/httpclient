package middleware

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestNewPrometheus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.Trim(r.URL.Path, "/")
		code, _ := strconv.Atoi(path)
		w.WriteHeader(code)
	}))
	defer srv.Close()

	numRequests := 10
	for i := 0; i < numRequests; i++ {
		middleware, _ := NewPrometheus()
		c := http.Client{Transport: WithMiddleware(nil, middleware)}
		if _, err := c.Get(srv.URL + "/200"); err != nil {
			t.Fatal(err)
		}
		if _, err := c.Get(srv.URL + "/500"); err != nil {
			t.Fatal(err)
		}
	}

	u, _ := url.Parse(srv.URL)
	assertCounter := func(t *testing.T, host, code string, expect int) {
		metrics, err := monitoring.counter.GetMetricWith(prometheus.Labels{
			"host": host, "path": "/200", "method": "GET", "response_status_code": "200",
		})
		if err == nil {
			if n := testutil.ToFloat64(metrics); n != float64(numRequests) {
				t.Errorf("Number of requests for %s code not equals. Expected %d, actual %f", code, expect, n)
			}
		} else {
			t.Errorf("Metrics has errors for %s code: %s", code, err)
		}
	}

	assertCounter(t, u.Host, "200", 10)
	assertCounter(t, u.Host, "500", 10)
}
