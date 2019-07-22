## Usage

You will need to configure opentracing.Tracer for your project.

Examples:
* [Zipkin](https://github.com/openzipkin/zipkin-go-opentracing/blob/master/examples/cli_with_2_services/cli/main.go)
* [Jaeger](https://github.com/opentracing-contrib/examples/blob/master/go/trivial.go)

You can configure opentracing.Tracer for your project and set it as a global tracer.

```go
opentrace.SetGlobalTracer(tracer)
````

Or you can set it for a specified transport.
```go
middleware.NewOpentrace().WithTracer(tracer)
```

#### Default
By default, the transport is using **StandardSpanner** and **GlobalTracer**.

**Tags**: span.kind, http.method, http.url, http.status_code, error.

```go
rt := http.DefaultTransport

client := &http.Client{
   Transport: middleware.NewOpentrace().
      RoundTripper(rt)
}

request, _ := http.NewRequest(http.MethodGet, "http://resource.io", nil)
request = request.WithContext(
   opentracing.ContextWithSpan(request.Context(), span),  
)

client.Do(request)

// For "StandardSpanner", you have to finish the span by yourself
span.Finish()
```

If you do not want to use GlobalTracer, you can set up a tracer directly:
```go
rt := http.DefaultTransport

client := &http.Client{
   Transport: middleware.NewOpentrace().
      WithTracer(tracer).
      RoundTripper(rt)
}
```

#### CreatorSpanner
Creates a span with a generated "operation name" and injects the tags.
**Tags**: span.kind, http.method, http.url, http.status_code, error.

```go
rt := http.DefaultTransport

client := &http.Client{
   Transport: middleware.NewOpentrace().
      WithSpanner(new(opentrace.CreatorSpanner))
      RoundTripper(rt)
   }


// Middleware will use the "root" span as a child for the new span
ctx := tracer.StartSpan("root")
r1, _ := http.NewRequest(http.MethodGet, "http://resource.io", nil)
r1 = request.WithContext(
   opentracing.ContextWithSpan(request.Context(), span),
)
client.Do(r1)

// We want to use the "root" span for this request, so we will skip the creation step
r2, _ := http.NewRequest(http.MethodGet, "http://resource.io", nil)
r2 = request.WithContext(ctx)
r2 = request.WithContext(
   opentrace.ContextWithSkipSpanCreating(request.Context()),
)
client.Do(r2)

// We did not set any span for this request
// By default, middleware does not create the "root" span
r3, _ := http.NewRequest(http.MethodGet, "http://resource.io", nil)
client.Do(r3)

// But you can enable this option
spanner := new(opentrace.CreatorSpanner).
   WithCreateRootSpanOnMissingParent(false)

&http.Client{
   Transport: middleware.NewOpentrace().
      WithSpanner(spanner).
      RoundTripper(rt)
   }
```

* For a RESTful URL, you can define a custom implementation of the naming function.
```go
spanner := opentrace.CreatorSpanner{
   OperationNameFn: func(r *http.Request) string {
      return "my-operation"
   }
}

client := &http.Client{
   Transport: middleware.NewOpentrace().
      WithSpanner(new(opentrace.CreatorSpanner))
      RoundTripper(rt)
}
```

#### Spanners
  
* **StandardSpanner** - Works only with an existing span from the Request context. will not work if a span does not exist.
* **CreatorSpanner** - Creates a new span and uses it as a child of the previous span.

#### Injector

* **HTTPHeadersInjector** - A wrapper over opentracing.HTTPHeadersCarrier that prevents modification of http.Request.