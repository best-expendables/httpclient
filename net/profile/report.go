package profile

import (
	"math"
	"time"
)

// Report for single HTTP Request
type Report struct {
	// ConnectStart time when client start connection
	ConnectStart time.Time
	// ConnectDone time when connection has been established
	ConnectDone time.Time

	// DNSLookupStart begin of DNS lookup
	DNSLookupStart time.Time
	// DNSLookupDone end of DNS lookup
	DNSLookupDone time.Time

	// TLSHandshakeStart begin of TLS handshake
	TLSHandshakeStart time.Time
	// TLSHandshakeDone end of TLS handshake
	TLSHandshakeDone time.Time

	// Reused connection from connection pool (keep-alive)
	Reused bool
}

// ConnectionTime time for establishing a connection
func (r *Report) ConnectionTime() time.Duration {
	return r.ConnectDone.Sub(r.ConnectStart)
}

// ConnectionTime time for establishing a connection in milliseconds
func (r *Report) ConnectionTimeMs() float64 {
	return toMilliseconds(r.ConnectionTime())
}

// DNSLookupTime time for DNS lookup
func (r *Report) DNSLookupTime() time.Duration {
	return r.DNSLookupDone.Sub(r.DNSLookupStart)
}

// DNSLookupTime time for DNS lookup in milliseconds
func (r *Report) DNSLookupTimeMs() float64 {
	return toMilliseconds(r.DNSLookupTime())
}

// TLSHandshakeDuration time for TLSHandshake
func (r *Report) TLSHandshakeTime() time.Duration {
	return r.TLSHandshakeDone.Sub(r.TLSHandshakeStart)
}

// TLSHandshakeTimeMs time for TLSHandshake in milliseconds
func (r *Report) TLSHandshakeTimeMs() float64 {
	return toMilliseconds(r.TLSHandshakeTime())
}

// toMilliseconds converts duration to milliseconds, precision 0.02.
func toMilliseconds(duration time.Duration) float64 {
	if duration < time.Microsecond*10 {
		return 0
	}

	ms := float64(duration) / float64(time.Millisecond)
	// Round time to 0.02 precision
	return math.Round(ms*100) / 100
}
