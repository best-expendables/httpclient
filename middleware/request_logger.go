package middleware

import (
	"net/http"

	log "github.com/best-expendables/logger"
)

// RequestLogger create log for request
type RequestLogger struct {
	logger
}

// NewRequestLogger create logger for request
func NewRequestLogger(loggerEntry log.Entry) *RequestLogger {
	return &RequestLogger{
		logger: logger{logger: loggerEntry},
	}
}

func (l *RequestLogger) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return RoundTripperFn(func(request *http.Request) (*http.Response, error) {
		if request != nil {
			l.Process(request)
		}

		return next.RoundTrip(request)
	})
}

func (l *RequestLogger) Process(request *http.Request) error {
	logger := l.logger.get(request.Context())

	if logger == nil {
		return nil
	}

	requestEntry, err := newRequestEntry(request)

	meta := log.Fields{
		"request": requestEntry,
		"source":  "RequestLogger",
		"url":     request.URL.String(),
	}

	if err != nil {
		meta["err"] = err
		logger.WithFields(meta).Warning("Request logger has an error")

		return nil
	}

	logger.WithFields(meta).Info()

	return nil
}
