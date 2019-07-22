package opentrace

import (
	"testing"

	"net/http"

	"errors"

	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/mocktracer"
)

func TestStandardSpanner_Open(t *testing.T) {
	tracer := mocktracer.New()

	request, _ := http.NewRequest("GET", "http://httptrace.io/resource", nil)
	request = request.WithContext(opentracing.ContextWithSpan(
		request.Context(),
		tracer.StartSpan("root"),
	))

	span := StandardSpanner{}.OnRequest(tracer, request).(*mocktracer.MockSpan)

	if operationName := span.OperationName; operationName != "root" {
		t.Errorf("operation name not equals: expected '%s', actual '%s'", "root", operationName)
	}

	// StandardSpanner inject 4 tags
	if count := len(span.Tags()); count != 3 {
		t.Errorf("number of tags not equals: expected '%d', actual '%d'", 4, count)
	}
}

func TestStandardSpanner_OnResponse(t *testing.T) {
	tracer := mocktracer.New()
	span := tracer.StartSpan("root").(*mocktracer.MockSpan)

	StandardSpanner{}.OnResponse(span, &http.Response{
		StatusCode: http.StatusOK,
	}, nil)

	if tag := span.Tag(string(ext.HTTPStatusCode)); tag == nil {
		t.Error("'HTTPStatusCode' is missing")
	} else {
		if statusCode, ok := tag.(uint16); !ok || statusCode != http.StatusOK {
			t.Errorf("status code not equals: expected '%d', actual '%d'", http.StatusOK, statusCode)
		}
	}
}

func TestStandardSpanner_OnResponse_ClientError(t *testing.T) {
	tracer := mocktracer.New()

	span := tracer.StartSpan("root").(*mocktracer.MockSpan)
	spanner := StandardSpanner{}
	spanner.OnResponse(span, nil, errors.New("client error"))

	if tag := span.Tag(string(ext.Error)); tag == nil {
		t.Error("tag 'error' is missing")
	}
}

func TestCreatorSpanner_Open_NewSpan_WithoutRootContext(t *testing.T) {
	tracer := mocktracer.New()

	request, _ := http.NewRequest("GET", "http://httptrace.io/resource", nil)

	spanner := CreatorSpanner{}

	// Request doesn't have span, span have be nil
	if span := spanner.OnRequest(tracer, request); span != nil {
		t.Error("span can not be nil")
	}

	// With this flag spanner should create span anyway
	spanner.WithCreateRootSpanOnMissingParent(true)
	span := spanner.OnRequest(tracer, request)

	if span == nil {
		t.Error("Span can not be nill")
	}

	if span, ok := span.(*mocktracer.MockSpan); ok {
		if span.OperationName != "/resource" {
			t.Errorf("operation name not equals: expected '%s', actual '%s'",
				"[GET] - http://opentrace.io/resource", span.OperationName)
		}
	} else {
		t.Error("Span is not *mocktracer.MockSpan")
	}

	// OnResponse should finish this span
	spanner.OnResponse(span, &http.Response{
		StatusCode: http.StatusOK,
		Request:    &http.Request{},
	}, nil)

	if len(tracer.FinishedSpans()) != 1 {
		t.Error("Span doesn't finished")
	}
}

func TestCreatorSpanner_Open_NewSpan_WithRootContext(t *testing.T) {
	tracer := mocktracer.New()

	request, _ := http.NewRequest("GET", "http://httptrace.io/resource", nil)
	request = request.WithContext(
		opentracing.ContextWithSpan(request.Context(), tracer.StartSpan("root")),
	)

	spanner := CreatorSpanner{
		OperationNameFn: func(r *http.Request) string {
			return "child"
		},
	}

	span := spanner.OnRequest(tracer, request)

	if span == nil {
		t.Error("span can not be nil")
	}

	if span, ok := span.(*mocktracer.MockSpan); ok {
		if span.OperationName != "child" {
			t.Errorf("operation name not equals: expected '%s', actual '%s'",
				"child", span.OperationName)
		}
	} else {
		t.Error("type assertion error")
	}
}

func TestCreatorSpanner_Open_WithSkipSpanCreating(t *testing.T) {
	tracer := mocktracer.New()

	request, _ := http.NewRequest("GET", "http://httptrace.io", nil)
	ctx := ContextWithSkipSpanCreating(opentracing.ContextWithSpan(
		context.Background(),
		tracer.StartSpan("root"),
	))
	request = request.WithContext(ctx)

	spanner := CreatorSpanner{}
	if span := spanner.OnRequest(tracer, request); span == nil {
		t.Error("span can not be nil")
	} else {
		span := span.(*mocktracer.MockSpan)

		if span.OperationName != "root" {
			t.Errorf("operation name not equals: expected '%s', actual '%s'",
				"root", span.OperationName)
		}
	}
}
