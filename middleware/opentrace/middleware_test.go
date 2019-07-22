package opentrace

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"errors"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestTransport_RoundTripper(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tracer := mocktracer.New()
	span := tracer.StartSpan("root")

	client := &http.Client{
		Transport: NewTransport().
			WithTracer(tracer).
			RoundTripper(http.DefaultTransport),
	}

	request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	request.Header.Set("Authorization", "Bearer")

	request = request.WithContext(
		opentracing.ContextWithSpan(request.Context(), span),
	)
	if _, err := client.Do(request); err != nil {
		t.Error(err)
	}

	span.Finish()

	// Headers have not been changed
	if len(request.Header) > 1 && request.Header.Get("Authorization") != "Bearer" {
		t.Errorf("request headers has been changed")
	}

	if count := len(tracer.FinishedSpans()); count != 1 {
		t.Errorf("number of finished spans '%d', expected 1", count)
	}

	finishedSpan := tracer.FinishedSpans()[0]

	if finishedSpan.OperationName != "root" {
		t.Errorf("operation name not equal: expected '%s', actual '%s'",
			"root", finishedSpan.OperationName)
	}

	if finishedSpan.FinishTime.IsZero() {
		t.Error("Finish time is zero")
	}
}

func TestTransport_RoundTripper_WithGlobalTracer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tracer := mocktracer.New()
	opentracing.SetGlobalTracer(tracer)
	span := tracer.StartSpan("root")

	client := &http.Client{
		Transport: NewTransport().RoundTripper(http.DefaultTransport),
	}

	request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	request = request.WithContext(
		opentracing.ContextWithSpan(request.Context(), span),
	)
	if _, err := client.Do(request); err != nil {
		t.Error(err)
	}

	span.Finish()

	if count := len(tracer.FinishedSpans()); count != 1 {
		t.Errorf("number of finished spans '%d', expected 1", count)
	}
}

func TestTransport_RoundTripper_SpannerReturnsNil(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tracer := mocktracer.New()

	client := &http.Client{
		Transport: NewTransport().
			WithTracer(tracer).
			RoundTripper(http.DefaultTransport),
	}

	request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	request.Header.Set("Authorization", "Bearer")
	if _, err := client.Do(request); err != nil {
		t.Error(err)
	}
}

func TestTransport_RoundTripper_InjectorError_InterruptOnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	tracer := mocktracer.New()
	span := tracer.StartSpan("root")

	client := &http.Client{
		Transport: NewTransport().
			WithTracer(tracer).
			WithInterruptOnError(true).
			WithInjector(InjectorFn(func(tracer opentracing.Tracer, ctx opentracing.SpanContext, r **http.Request) error {
				return errors.New("internal error")
			})).
			RoundTripper(http.DefaultTransport),
	}

	request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
	request = request.WithContext(
		opentracing.ContextWithSpan(request.Context(), span),
	)

	if _, err := client.Do(request); err == nil {
		t.Error("error can not be nil")
	}
}

func TestTransport_WithSpanner(t *testing.T) {
	transport := NewTransport()

	spanner := transport.spanner
	spanner2 := transport.WithSpanner(new(CreatorSpanner)).spanner

	if spanner == spanner2 {
		t.Error("Span hasn't changed")
	}
}
