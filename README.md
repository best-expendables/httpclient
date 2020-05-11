# HTTP Client GO library

### Install
```yaml
import:
- package: github.com/best-expendables/httpclient
  version: x.x.x
```

### Middlewares

- `middleware.Authentication`
- `middleware.RequestLogger`
- `middleware.ResponseLogger`
- `middleware.Newrelic`
- `middleware.NewNewrelicApiGateway`
-  [middleware.Opentrace](https://bitbucket.lzd.co/projects/LGO/repos/httpclient/browse/docs/opentrace.md)
- `middleware.NetworkProfiler`
- `middleware.RequestID`

#### RequestLogger/ResponseLogger
Since we use request-dependent logging, we have to pass context with logger to each request.  
Otherwise middleware will use logger which was injected via constructor `middleware.NewRequestLogger(logger)` 

#### NetworkProfiler
Network profiler collects metrics about the network and set the report into context.  
Low overhead cost allows to use it for production.  

 
You can find the report into the response log:
```text
{"content":{"response":{...}, "network": {...}}

"reused": (bool) - Connection was taken from the keep-alive pool
"connection": Time elapsed for connection establishing (ignores on keep-alive connection)
"dns": Time elapsed for DNS lookup (ignores on keep-alive connection)

```


Or for some cases you can get the report from response e.g.:
```go
package main(
	"net/http"
    
	"github.com/best-expendables/httpclient/middleware"
	"github.com/best-expendables/httpclient/net/profile"
)
	

transport = middleware.WithMiddleware(nil, NewNetworkProfiler)

c := &http.Client{
	Transport: transport,
}
response, _ := c.Get("http://localhost")
report := profile.ReportFromResponse(response)
````

### Examples


##### 1) Context-independent middlewares. 
One client for whole life time of application
We can keep one client for whole application lifetime and pass context for each http request.

```go
package main

import (
	"net/http"

	"github.com/best-expendables/httpclient/middleware"
	log "github.com/best-expendables/logger"
)

func main() {
	logger := log.NewLoggerFactory(log.InfoLevel).
		Logger(context.TODO())

	transport = middleware.WithMiddleware(
		http.DefaultTransport,
		middleware.NewRequestLogger(logger),
		middleware.NewResponseLogger(logger),
		middleware.NewOpentrace(),
	)
	
	client := http.Client{
		Transport: transport,
	}
	
	request, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	request = request.WithContext(ctx)
}
```

##### 2) Context-dependent middlewares. 
New client for each incoming request.   
For compatibility with old code, we can create middleware for each http incoming request.   
In this case application will spend more resources for memory allocation and garbage collecting.   

```go
package main

import (
	"net/http"

	"github.com/best-expendables/httpclient/middleware"
	log "github.com/best-expendables/logger"
)

func main() {
	logger := log.NewLoggerFactory(log.InfoLevel).
		Logger(context.TODO())

	transport = middleware.WithMiddleware(
		http.DefaultTransport,
		middleware.NewRequestLogger(logger),
		middleware.NewResponseLogger(logger),
		middleware.NewOpentrace(),
	)
	
	client := http.Client{
		Transport: transport,
	}
	
	request, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
	request = request.WithContext(ctx)
}
```