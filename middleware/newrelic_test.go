package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type NewRelicSuite struct {
	suite.Suite
}

func (s *NewRelicSuite) Test() {
	client := &http.Client{
		Transport: WithMiddleware(nil, NewNewrelic(nil)),
	}

	s.T().Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		server.URL += "/user/1000/details"

		s.T().Run("GET", func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
			response, err := client.Do(request)
			s.NoError(err)
			s.NotNil(response)
		})

		s.T().Run("Nil context", func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, server.URL, nil)
			response, err := client.Do(request)
			s.NoError(err)
			s.NotNil(response)
		})
	})

	s.T().Run("ApiGateway formatter", func(t *testing.T) {
		nr := NewNewrelicApiGateway()
		r, _ := http.NewRequest(http.MethodGet, "http://apigateway/v1/addresses/region/11321", nil)
		expected := "http://apigateway.v1.addresses"
		if actual := nr.formatter.ToURL(r); actual != expected {
			t.Errorf("URL not equals. Expected '%s', actual '%s'", expected, actual)
		}
	})
}

func TestNewRelicRunner(t *testing.T) {
	suite.Run(t, new(NewRelicSuite))
}
