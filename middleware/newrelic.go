package middleware

import (
	"net/http"
	"strings"

	log "bitbucket.org/snapmartinc/logger"
	newrelic "github.com/newrelic/go-agent"
)

type (
	// NewrelicRequestFormatter custom URL parsing
	NewrelicRequestFormatter interface {
		ToURL(request *http.Request) string
	}

	// Newrelic middleware for HTTP external request
	Newrelic struct {
		formatter NewrelicRequestFormatter
	}

	// NewrelicRequestFormatterFunc implements NewrelicRequestFormatter interface
	NewrelicRequestFormatterFunc func(r *http.Request) string
)

// NewNewrelic formatter could be nil, than URL will be parsed via newrelic
func NewNewrelic(formatter NewrelicRequestFormatter) *Newrelic {
	return &Newrelic{formatter: formatter}
}

// NewNewrelicApiGateway should be used for the calls through API Gateway
func NewNewrelicApiGateway() *Newrelic {
	return NewNewrelic(NewrelicRequestFormatterFunc(func(r *http.Request) string {
		path := strings.Trim(r.URL.Path, "/")
		paths := strings.SplitN(path, "/", 3)
		url := r.URL.Scheme + "://" + r.URL.Host
		// It's 2 because we always have /v1 or /v2 as a first segment of URL path
		if len(paths) >= 2 {
			url = url + "." + paths[0] + "." + paths[1]
		}
		return url
	}))
}

func (r Newrelic) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFn(func(request *http.Request) (*http.Response, error) {
		txn := newrelic.FromContext(request.Context())
		if txn != nil {
			segment := newrelic.StartExternalSegment(txn, request)
			if r.formatter != nil {
				segment.URL = r.formatter.ToURL(request)
			}
			defer func() {
				if err := segment.End(); err != nil {
					log.EntryFromContextOrDefault(request.Context()).WithFields(log.Fields{
						"err":       err.Error(),
						"component": "httpclient.newrelic",
					}).Error("Can not finish newrelic segment")
				}
			}()
		}
		return next.RoundTrip(request)
	})
}

func (fn NewrelicRequestFormatterFunc) ToURL(r *http.Request) string {
	return fn(r)
}
