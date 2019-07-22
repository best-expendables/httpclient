package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"strings"

	log "bitbucket.org/snapmartinc/logger"
	"github.com/stretchr/testify/suite"
)

type ResponseLoggerSuite struct {
	logger      log.Entry
	buffer      *bytes.Buffer
	requestBody string

	suite.Suite
}

func (s *ResponseLoggerSuite) Test_LoggerFromProperty() {
	client := &http.Client{
		Transport: WithMiddleware(nil, NewResponseLogger(s.logger)),
	}

	s.T().Run("Success", func(t *testing.T) {
		responseBody := `{ "status": "OK" }`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(responseBody))
		}))
		defer server.Close()

		s.T().Run("GET", func(t *testing.T) {
			response, err := client.Get(server.URL)
			s.NoError(err)
			s.NotNil(response)
			s.LogNotEmpty()

			t.Log(s.buffer.String())
			s.buffer.Reset()
		})

		s.T().Run("POST", func(t *testing.T) {
			response, err := client.Post(server.URL, "", bytes.NewBufferString(s.requestBody))
			s.NoError(err)
			s.NotNil(response)
			s.LogNotEmpty()

			t.Log(s.buffer.String())
			s.buffer.Reset()
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

			t.Log(s.buffer.String())
			s.buffer.Reset()
		})
	})

	s.T().Run("Validation error response code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusUnprocessableEntity)
		}))
		defer server.Close()

		response, err := client.Get(server.URL)
		s.NoError(err)
		s.NotNil(response)
		s.LogNotEmpty()

		responseLog := s.buffer.String()
		t.Log(responseLog)
		if !strings.Contains(responseLog, "\"level\":\"info\"") {
			t.Error("Should be info level")
		}

		if !strings.Contains(responseLog, http.StatusText(http.StatusUnprocessableEntity)) {
			t.Error()
		}

		s.buffer.Reset()
	})
}

func (s *ResponseLoggerSuite) Test_LoggerFromContext() {
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
		t.Log(s.buffer.String())
		s.buffer.Reset()
	})
}

func (s *ResponseLoggerSuite) Test_LoggerWithProfiler() {
	client := &http.Client{
		Transport: WithMiddleware(nil, NewResponseLogger(nil), NewNetworkProfiler()),
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

		logRecord := s.buffer.String()

		if !strings.Contains(logRecord, "\"network\"") {
			t.Error("Log doesn't contains network profiling information")
		}

		t.Log(logRecord)
		s.buffer.Reset()
	})
}

func (s *ResponseLoggerSuite) LogNotEmpty() bool {
	return s.NotEmpty(s.buffer.String())
}

func TestLoggerRunner(t *testing.T) {
	buffer := bytes.NewBufferString("")
	logger := log.NewLoggerFactory(log.InfoLevel, log.SetOut(buffer)).Logger(
		context.TODO(),
	)

	suite.Run(t, &ResponseLoggerSuite{
		logger:      logger,
		buffer:      buffer,
		requestBody: `{ "name":"John", "age":30, "car":null }`,
	})
}
