package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AuthorizationSuite struct {
	suite.Suite
}

func (s *AuthorizationSuite) Test() {
	middleware := NewAuthentication("0aac4e6a54c170b06e2bd3848d2b735e")

	client := &http.Client{
		Transport: WithMiddleware(nil, middleware),
	}

	s.T().Run("Success", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			headerName, headerValue := middleware.Header()
			header := r.Header.Get(headerName)
			s.Equal(header, headerValue)
			w.WriteHeader(http.StatusOK)
		}))
		defer server.Close()

		s.T().Run("GET", func(t *testing.T) {
			response, err := client.Get(server.URL)
			s.NoError(err)
			s.NotNil(response)
		})

		s.T().Run("POST", func(t *testing.T) {
			response, err := client.Post(server.URL, "", nil)
			s.NoError(err)
			s.NotNil(response)
		})
	})
}

func TestAuthorizationRunner(t *testing.T) {
	suite.Run(t, new(AuthorizationSuite))
}
