package middleware

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

var monitoring *Prometheus

type Prometheus struct {
	counter *prometheus.CounterVec
}

// NewPrometheus
func NewPrometheus() (*Prometheus, error) {
	if monitoring != nil {
		return monitoring, nil
	}

	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: "http_outgoing_requests",
		Name:      "total",
	}, []string{"host", "path", "method", "response_status_code"})

	if err := prometheus.Register(counter); err != nil {
		return nil, err
	}

	monitoring = &Prometheus{counter: counter}

	return monitoring, nil
}

func (p *Prometheus) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFn(func(request *http.Request) (*http.Response, error) {
		response, err := next.RoundTrip(request)

		if err != nil {
			return response, err
		}

		p.counter.With(prometheus.Labels{
			"host":               request.URL.Host,
			"path":               request.URL.Path,
			"method":             request.Method,
			"response_status_code": strconv.Itoa(response.StatusCode),
		}).Inc()

		return response, err
	})
}
