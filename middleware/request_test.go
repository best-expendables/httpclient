package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRequestEntry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	request, _ := http.NewRequest(http.MethodGet, server.URL, bytes.NewBufferString("test body"))
	request.Header.Set("Test", "1")

	requestEntry, err := newRequestEntry(request)

	read, _ := ioutil.ReadAll(request.Body)

	if len(read) == 0 {
		t.Error("Request body should be recovered")
	}

	if requestEntry.Body != "test body" {
		t.Error(fmt.Sprintf("expect %s to be \"test body\"", request.Body))
	}

	if requestEntry.Method != http.MethodGet {
		t.Error("expect method to be GET but got", requestEntry.Method)
	}

	if requestEntry.binaryBody {
		t.Error("expect binaryBody to be false but true")
	}

	if err != nil {
		t.Error("expect no errors")
	}

	if val := requestEntry.Header.Get("Test"); val != "1" {
		t.Error("expect value to be 1 but got", val)
	}
}

func TestRestEntry_parseFormToJson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	values := url.Values{}

	values.Add("key1", "value1")
	values.Add("key2", "value2")
	body := values.Encode()

	request, _ := http.NewRequest(http.MethodPost, server.URL, strings.NewReader(body))

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	parsedBody := parseFormToJson(request, bytes.NewBufferString(body).Bytes())

	if parsedBody != "{\"key1\":[\"value1\"],\"key2\":[\"value2\"]}" {
		t.Error("expected body to be parsed correctly")
	}
}
