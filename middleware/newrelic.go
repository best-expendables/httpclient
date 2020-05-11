package middleware

import (
	"net/http"
	"strings"

	log "github.com/best-expendables/logger"
	newrelic "github.com/newrelic/go-agent"
)

type URLFormatFunc func(r *http.Request) string

type Newrelic struct {
	urlFormatter URLFormatFunc
}

func NewURLFormatFunc() URLFormatFunc {
	return func(r *http.Request) string {
		path := strings.Trim(r.URL.Path, "/")
		paths := strings.SplitN(path, "/", 3)
		url := r.URL.Scheme + "://" + r.URL.Host
		// It's 2 because we always have /v1 or /v2 as a first segment of URL path
		if len(paths) >= 2 {
			url = url + "." + paths[0] + "." + paths[1]
		}
		return url
	}
}

func NewNewrelicApiGateway(urlFormatter URLFormatFunc) *Newrelic {
	return &Newrelic{urlFormatter: urlFormatter}
}

func (r Newrelic) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFn(func(request *http.Request) (*http.Response, error) {
		txn := newrelic.FromContext(request.Context())
		if txn != nil {
			segment := newrelic.StartExternalSegment(txn, request)
			if r.urlFormatter != nil {
				segment.URL = r.urlFormatter(request)
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
