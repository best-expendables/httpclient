package opentrace

import (
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type (
	// Spanner injects tags, logs and manages the lifecycle of a span
	//
	// Tagging and logging specification:
	// https://github.com/opentracing/specification/blob/master/semantic_conventions.md
	Spanner interface {
		// OnRequest called on an HTTP request
		OnRequest(tracer opentracing.Tracer, request *http.Request) opentracing.Span

		// OnResponse called on an HTTP response
		OnResponse(span opentracing.Span, response *http.Response, clientError error)
	}

	// StandardSpanner is used by default
	// It only works with the existing span and does not create a new one
	StandardSpanner struct{}

	// CreatorSpanner - extended standard spanner with one difference - it creates a new span on each request
	CreatorSpanner struct {
		// CreateRootSpanOnMissingParent creates a "root" span if the parent span is missing
		CreateRootSpanOnMissingParent bool

		// OperationNameFn function creates an operation name for the new span
		OperationNameFn func(r *http.Request) string

		StandardSpanner
	}
)

// OnRequest adds tags to an existing span: span.kind, http.url, http.method
func (StandardSpanner) OnRequest(tracer opentracing.Tracer, request *http.Request) opentracing.Span {
	var span opentracing.Span

	if span = opentracing.SpanFromContext(request.Context()); span == nil {
		return nil
	}

	ext.SpanKindRPCClient.Set(span)

	ext.HTTPUrl.Set(span, request.URL.String())
	ext.HTTPMethod.Set(span, request.Method)

	return span
}

// OnResponse adds a "HTTPStatusCode" or "Error" tag to the existing span
func (StandardSpanner) OnResponse(span opentracing.Span, response *http.Response, clientError error) {
	if clientError != nil {
		ext.Error.Set(span, true)
	}

	if response != nil {
		ext.HTTPStatusCode.Set(span, uint16(response.StatusCode))
	}
}

// OnRequest creates a new span and passes it to StandardSpanner
func (p *CreatorSpanner) OnRequest(tracer opentracing.Tracer, request *http.Request) opentracing.Span {
	if SkipSpanCreatingFromContext(request.Context()) {
		return p.StandardSpanner.OnRequest(tracer, request)
	}

	var (
		span          opentracing.Span
		operationName string
	)

	if p.OperationNameFn == nil {
		operationName = OperationNameFromRequest(request)
	} else {
		operationName = p.OperationNameFn(request)
	}

	if parentSpan := opentracing.SpanFromContext(request.Context()); parentSpan != nil {
		span = tracer.StartSpan(
			operationName, opentracing.ChildOf(parentSpan.Context()),
		)
	} else if p.CreateRootSpanOnMissingParent {
		span = tracer.StartSpan(operationName)
	} else {
		return nil
	}

	ctx := opentracing.ContextWithSpan(request.Context(), span)

	return p.StandardSpanner.OnRequest(tracer, request.WithContext(ctx))
}

// OnResponse passes a span to StandardSpanner
func (p *CreatorSpanner) OnResponse(span opentracing.Span, response *http.Response, clientError error) {
	p.StandardSpanner.OnResponse(span, response, clientError)

	if !SkipSpanCreatingFromContext(response.Request.Context()) {
		span.Finish()
	}
}

// WithCreateRootSpanOnMissingParent creates a "root" span if the parent span is missing
func (p *CreatorSpanner) WithCreateRootSpanOnMissingParent(flag bool) *CreatorSpanner {
	p.CreateRootSpanOnMissingParent = flag
	return p
}

// OperationNameFromRequest creates an operation name from the request
//
// E.g.:
// 	-	out: [POST] /v1/users
//	-	out: [GET] /v1/users
//
// For RESTful, you will need to use your own implementation
func OperationNameFromRequest(request *http.Request) string {
	return request.URL.Path
}
