package middleware

import (
	log "github.com/best-expendables/logger"
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
)

// loggable structure helper
type logger struct {
	logger log.Entry
}

// get context-dependent logger.
// If logger not presented into context then returns "base" logger from property.
func (l *logger) get(ctx context.Context) log.Entry {
	if logger := log.EntryFromContext(ctx); logger != nil {
		return logger
	}

	return l.logger
}

func cloneRequest(r *http.Request) *http.Request {
	r2 := new(http.Request)
	*r2 = *r
	// deep copy of the Header
	r2.Header = make(http.Header)
	for k, s := range r.Header {
		r2.Header[k] = s
	}
	return r2
}

func recoverRequestBody(request *http.Request, originalBody []byte) {
	if originalBody != nil {
		request.Body = ioutil.NopCloser(bytes.NewBuffer(originalBody))
	}
}
