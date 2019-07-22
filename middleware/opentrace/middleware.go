package opentrace

import (
	"github.com/opentracing/opentracing-go"
	"net/http"
)

type (
	// Transport is used for the HTTP Client transport, it implements the RoundTripper interface
	Transport struct {
		// InterruptOnError
		// E.g. HTTP request will be interrupted upon a tracer.Inject error
		InterruptOnError bool

		spanner  Spanner
		injector Injector
		tracer   opentracing.Tracer
	}

	roundTripperFn func(r *http.Request) (*http.Response, error)
)

// NewTransport returns http.RoundTripper
func NewTransport() *Transport {
	return &Transport{
		spanner:  new(StandardSpanner),
		injector: new(HTTPHeadersInjector),
	}
}

// Transport middleware function
func (o *Transport) RoundTripper(next http.RoundTripper) http.RoundTripper {
	return roundTripperFn(func(request *http.Request) (*http.Response, error) {
		var (
			span   opentracing.Span
			tracer opentracing.Tracer
		)

		if tracer = o.tracer; tracer == nil {
			tracer = opentracing.GlobalTracer()
		}

		if span = o.spanner.OnRequest(tracer, request); span == nil {
			return next.RoundTrip(request)
		}

		r := request
		if err := o.injector.Inject(tracer, span.Context(), &r); err != nil && o.InterruptOnError {
			return nil, err
		}

		response, err := next.RoundTrip(r)
		o.spanner.OnResponse(span, response, err)

		return response, err
	})
}

// WithInterruptOnError sets a flag
func (o *Transport) WithInterruptOnError(flag bool) *Transport {
	o.InterruptOnError = flag
	return o
}

// WithSpanner sets a spanner
func (o *Transport) WithSpanner(spanner Spanner) *Transport {
	o.spanner = spanner
	return o
}

// WithInjector sets an injector
func (o *Transport) WithInjector(injector Injector) *Transport {
	o.injector = injector
	return o
}

// WithTracer sets a tracer
func (o *Transport) WithTracer(tracer opentracing.Tracer) *Transport {
	o.tracer = tracer
	return o
}

// RoundTrip
func (f roundTripperFn) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
