package main

import (
	"net/http"

	"context"
	"net/http/httptest"

	"bitbucket.org/snapmartinc/httpclient/middleware"
	log "bitbucket.org/snapmartinc/logger"
	"net/url"
	"strings"
)

func main() {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// We can keep one client for whole application lifetime and pass context for each http request.
	{
		logger := log.NewLoggerFactory(log.InfoLevel).
			Logger(context.TODO())

		transport := middleware.WithMiddleware(
			http.DefaultTransport,
			middleware.NewRequestLogger(logger),
			middleware.NewResponseLogger(logger),
			middleware.NewOpentrace(),
		)

		client := http.Client{
			Transport: transport,
		}

		for i := 0; i < 3; i++ {
			// Some context from router
			ctx := context.Background()

			request, _ := http.NewRequest(http.MethodGet, srv.URL, nil)
			request.Header.Set("X-LEL-User-ID", "4ac32f8f-e746-47f0-8e1b-898ae9f5b80c")

			request = request.WithContext(ctx)

			client.Do(request)
		}

		body := url.Values{}

		body.Add("user", "user1")
		body.Add("pass", "pass1")

		for i := 0; i < 2; i++ {
			postRequest, _ := http.NewRequest(http.MethodPost, srv.URL, strings.NewReader(body.Encode()))

			if i == 1 {
				postRequest.Header.Add("content-type", "application/x-www-form-urlencoded")
			}

			client.Do(postRequest)
		}
	}

	// For compatibility with old code, we can create middleware for each http incoming request.
	//
	// In this case application will spend more resources for memory allocation and garbage collecting.
	{
		loggerFactory := log.NewLoggerFactory(log.InfoLevel)

		for i := 0; i < 3; i++ {
			// Some context from router
			ctx := context.Background()

			logger := loggerFactory.Logger(ctx)

			transport := middleware.WithMiddleware(
				http.DefaultTransport,
				middleware.NewRequestLogger(logger),
				middleware.NewResponseLogger(logger),
				middleware.NewOpentrace(),
			)

			request, _ := http.NewRequest(http.MethodGet, srv.URL, nil)

			client := http.Client{
				Transport: transport,
			}

			client.Do(request)
		}
	}
}
