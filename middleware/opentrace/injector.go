package opentrace

import (
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type (
	// Injector knows how to propagate a trace
	// http.Request agrument should be immutable
	// For any modifications, use a copy of the http.Request and replace the pointer
	Injector interface {
		Inject(tracer opentracing.Tracer, ctx opentracing.SpanContext, r **http.Request) error
	}

	// InjectorFn a wrapper for the Injector interface
	InjectorFn func(tracer opentracing.Tracer, ctx opentracing.SpanContext, r **http.Request) error

	// HTTPHeadersInjector used by default, it injects a trace into the HTTP headers
	HTTPHeadersInjector struct{}
)

// Inject
func (fn InjectorFn) Inject(tracer opentracing.Tracer, ctx opentracing.SpanContext, r **http.Request) error {
	return fn(tracer, ctx, r)
}

// Inject creates a copy of the http.Request and replaces a pointer in the argument
func (HTTPHeadersInjector) Inject(tracer opentracing.Tracer, ctx opentracing.SpanContext, r **http.Request) error {
	request := **r
	header := make(http.Header)
	for k, v := range request.Header {
		header[k] = v
	}
	request.Header = header

	err := tracer.Inject(ctx, opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(request.Header))
	*r = &request

	return err
}
