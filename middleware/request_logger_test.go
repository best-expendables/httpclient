package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"net/url"
	"strings"

	log "bitbucket.lzd.co/lgo/logger"
	"bitbucket.lzd.co/lgo/logger/mock"
	"github.com/stretchr/testify/suite"
)

type RequestLoggerSuite struct {
	logger      log.Entry
	buffer      *bytes.Buffer
	requestBody string

	suite.Suite
}

func (s *RequestLoggerSuite) Test_LoggerFromProperty() {
	client := &http.Client{
		Transport: WithMiddleware(nil, NewRequestLogger(s.logger)),
	}

	s.T().Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		s.T().Run("GET", func(t *testing.T) {
			response, err := client.Get(server.URL)
			s.NoError(err)
			s.NotNil(response)
			s.LogNotEmpty()

			s.buffer.Reset()
		})

		s.T().Run("POST", func(t *testing.T) {
			response, err := client.Post(server.URL, "", bytes.NewBufferString(s.requestBody))
			s.NoError(err)
			s.NotNil(response)
			s.LogNotEmpty()

			s.buffer.Reset()
		})

		s.T().Run("POST", func(t *testing.T) {
			data := url.Values{}

			data.Add("user", "1")
			data.Add("pass", "2")

			response, err := client.Post(server.URL, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))

			s.NoError(err)
			s.NotNil(response)
			s.LogNotEmpty()
			s.LogContains("\\\"user\\\":[\\\"1\\\"]")
			s.LogContains("\\\"pass\\\":[\\\"2\\\"]")
			s.buffer.Reset()
		})

		s.T().Run("GET. Logger from context", func(t *testing.T) {
			client := &http.Client{
				Transport: WithMiddleware(nil, NewRequestLogger(s.logger)),
			}

			calledFromCTX := false
			ctx := log.ContextWithEntry(&logmock.Entry{
				WithFieldsFn: func(fields log.Fields) log.Entry {
					return &logmock.Entry{
						InfoFn: func(args ...interface{}) {
							calledFromCTX = true
						},
					}
				},
			}, context.Background())

			request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
			request = request.WithContext(ctx)

			client.Do(request)

			s.buffer.Reset()
			s.True(calledFromCTX)
		})
	})

	s.T().Run("HTTP error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		s.T().Run("GET", func(t *testing.T) {
			response, err := client.Get(server.URL)
			s.NoError(err)
			s.NotNil(response)
			s.LogNotEmpty()

			s.buffer.Reset()
		})

		s.T().Run("POST", func(t *testing.T) {
			response, err := client.Post(server.URL, "", bytes.NewBufferString(s.requestBody))
			s.NoError(err)
			s.NotNil(response)
			s.LogNotEmpty()

			s.buffer.Reset()
		})
	})
}

func (s *RequestLoggerSuite) Test_FromContext() {
	client := &http.Client{
		Transport: WithMiddleware(nil, NewResponseLogger(nil)),
	}

	s.T().Run("Success", func(t *testing.T) {
		responseBody := `{ "status": "OK" }`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(responseBody))
		}))
		defer server.Close()

		logger := log.NewLoggerFactory(log.InfoLevel, log.SetOut(s.buffer)).Logger(context.TODO())

		request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
		request = request.WithContext(
			log.ContextWithEntry(logger, context.TODO()),
		)

		response, err := client.Do(request)

		s.NoError(err)
		s.NotNil(response)
		s.LogNotEmpty()

		s.buffer.Reset()
	})
}

func (s *RequestLoggerSuite) LogNotEmpty() bool {
	return s.NotEmpty(s.buffer.String())
}

func (s *RequestLoggerSuite) RequestBodyNotEmpty() bool {
	return s.NotEmpty(s.requestBody)
}

func (s *RequestLoggerSuite) LogContains(substring string) bool {
	test := s.buffer.String()
	return s.True(strings.Contains(test, substring))
}

func TestRequestLoggerSuite(t *testing.T) {
	buffer := bytes.NewBufferString("")
	logger := log.NewLoggerFactory(log.InfoLevel, log.SetOut(buffer)).Logger(context.TODO())

	suite.Run(t, &RequestLoggerSuite{
		logger:      logger,
		buffer:      buffer,
		requestBody: `{ "name":"John", "age":30, "car":null }`,
	})
}
