package httpclient_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"bitbucket.org/snapmartinc/httpclient"
	"github.com/stretchr/testify/assert"
)

type CustomResponseParser struct{}

func (c *CustomResponseParser) Parse(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

func TestBaseClient_DoRequest(t *testing.T) {
	mockHandler := func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "application/json", req.Header.Get("content-type"))
		param := req.URL.Query().Get("query")
		if param == "not-found" {
			rw.WriteHeader(http.StatusInternalServerError)
			msg := `internal server error`
			rw.Write([]byte(msg))
			return
		}
		data := `{
			"data": {
				"field_1": "field 1 data",
				"field_2": 100
			    }
		}`
		rw.Write([]byte(data))
	}

	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()
	c := httpclient.NewBaseClient(server.URL)

	params := url.Values{
		"query": []string{"test"},
	}
	result := struct {
		Field1 string `json:"field_1"`
		Field2 int    `json:"field_2"`
	}{}
	err := c.DoRequest(context.Background(), "GET", "/test-path", params, nil, &result)
	assert.Nil(t, err)
	assert.Equal(t, "field 1 data", result.Field1)
	assert.Equal(t, 100, result.Field2)

	params = url.Values{
		"query": []string{"not-found"},
	}
	err = c.DoRequest(context.Background(), "GET", "/test-path", params, nil, nil)
	assert.NotNil(t, err)
	assert.IsType(t, &httpclient.Error{}, err)
	concreteErr, _ := err.(*httpclient.Error)
	assert.Equal(t, "internal server error", concreteErr.Message)
}

func TestBaseClient_WithResponseParser(t *testing.T) {

	mockHandler := func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "application/json", req.Header.Get("content-type"))
		param := req.URL.Query().Get("query")
		if param == "not-found" {
			rw.WriteHeader(http.StatusInternalServerError)
			msg := `internal server error`
			rw.Write([]byte(msg))
			return
		}
		data := `{
			"field_1": "field 1 data",
			"field_2": 100
		}`
		rw.Write([]byte(data))
	}

	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()

	c := httpclient.NewBaseClient(server.URL, httpclient.WithResponseParser(&CustomResponseParser{}))

	params := url.Values{
		"query": []string{"test"},
	}
	result := struct {
		Field1 string `json:"field_1"`
		Field2 int    `json:"field_2"`
	}{}
	err := c.DoRequest(context.Background(), "GET", "/test-path", params, nil, &result)
	assert.Nil(t, err)
	assert.Equal(t, "field 1 data", result.Field1)
	assert.Equal(t, 100, result.Field2)

	params = url.Values{
		"query": []string{"not-found"},
	}
	err = c.DoRequest(context.Background(), "GET", "/test-path", params, nil, nil)
	assert.NotNil(t, err)
	assert.IsType(t, &httpclient.Error{}, err)
	concreteErr, _ := err.(*httpclient.Error)
	assert.Equal(t, "internal server error", concreteErr.Message)
}
