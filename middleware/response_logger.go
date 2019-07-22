package middleware

import (
	"httpclient/net/profile"
	"net/http"

	log "bitbucket.org/snapmartinc/logger"
)

// ResponseLogger create log for response
type ResponseLogger struct {
	logger
}

// NewResponseLogger create logger for response
func NewResponseLogger(loggerEntry log.Entry) *ResponseLogger {
	return &ResponseLogger{
		logger: logger{logger: loggerEntry},
	}
}

func (l *ResponseLogger) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFn(func(request *http.Request) (*http.Response, error) {
		response, err := next.RoundTrip(request)
		if response != nil {
			l.Process(response)
		}

		return response, err
	})
}

func (l *ResponseLogger) Process(response *http.Response) error {
	logger := l.logger.get(response.Request.Context())

	if logger == nil {
		return nil
	}

	responseEntry, err := newResponseEntry(response)

	meta := log.Fields{
		"url":      response.Request.URL.String(),
		"source":   "ResponseLogger",
		"response": responseEntry,
	}

	if report := profile.ReportFromResponse(response); report != nil {
		network := make(map[string]interface{})
		network["reused"] = report.Reused

		if !report.Reused {
			network["dns"] = report.ConnectionTimeMs()
			network["connection"] = report.ConnectionTimeMs()
		}

		meta["network"] = network
	}

	if err != nil {
		meta["err"] = err
		logger.WithFields(meta).Warning("Response logger has an error")

		return nil
	}

	entry := logger.WithFields(meta)

	statusText := http.StatusText(response.StatusCode)

	// 422 response code uses for validation error
	if response.StatusCode < 400 || response.StatusCode == http.StatusUnprocessableEntity {
		entry.Info(statusText)
	} else {
		entry.Error(statusText)
	}

	return nil
}
