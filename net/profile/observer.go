package profile

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptrace"
	"time"
)

type ctxKey int

const reportCtxKey ctxKey = iota

// Observe request
func Observe(r *http.Request) *http.Request {
	report := new(Report)

	ctx := httptrace.WithClientTrace(
		r.Context(), observer(report),
	)

	ctx = context.WithValue(ctx, reportCtxKey, report)

	return r.WithContext(ctx)
}

// ReportFromResponse return report from response
func ReportFromResponse(response *http.Response) *Report {
	ctx := response.Request.Context()

	if report, ok := ctx.Value(reportCtxKey).(*Report); ok {
		return report
	}

	return nil
}

func observer(report *Report) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GotConn: func(conn httptrace.GotConnInfo) {
			report.Reused = conn.Reused
		},
		ConnectStart: func(_, _ string) {
			report.ConnectStart = time.Now()
		},
		ConnectDone: func(_, _ string, err error) {
			if err == nil {
				report.ConnectDone = time.Now()
			}
		},
		DNSStart: func(_ httptrace.DNSStartInfo) {
			report.DNSLookupStart = time.Now()
		},
		DNSDone: func(_ httptrace.DNSDoneInfo) {
			report.DNSLookupDone = time.Now()
		},
		TLSHandshakeStart: func() {
			report.TLSHandshakeStart = time.Now()
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, err error) {
			if err == nil {
				report.TLSHandshakeDone = time.Now()
			}
		},
	}
}
