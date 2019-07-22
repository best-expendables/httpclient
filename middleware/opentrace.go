package middleware

import "bitbucket.lzd.co/lgo/httpclient/middleware/opentrace"

// NewOpentrace create middleware for OpenTrace
func NewOpentrace() *opentrace.Transport {
	return opentrace.NewTransport()
}
