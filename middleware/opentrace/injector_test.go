package opentrace

import (
	"net/http"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestHTTPHeadersInjector_Inject(t *testing.T) {
	tracer := mocktracer.New()

	span := tracer.StartSpan("root")

	immutableRequest, _ := http.NewRequest(http.MethodGet, "http://httptrace.io/resource", nil)
	immutableRequest = immutableRequest.WithContext(opentracing.ContextWithSpan(immutableRequest.Context(), span))
	immutableRequest.Header.Set("HEADER", "HEADER")

	mutableRequest := immutableRequest
	injector := HTTPHeadersInjector{}
	if err := injector.Inject(tracer, span.Context(), &mutableRequest); err != nil {
		t.Error(err)
	}

	if len(immutableRequest.Header) > 1 || immutableRequest.Header.Get("HEADER") != "HEADER" {
		t.Error("Immutable request has been changed")
	}

	if mutableRequest.Header.Get("HEADER") != "HEADER" {
		t.Error("Muttable request lost header")
	}

	if len(mutableRequest.Header) == 1 {
		t.Error("Injector doesn't injected any header")
	}
}
