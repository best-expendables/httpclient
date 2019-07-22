package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponseEntry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	request, _ := http.NewRequest(http.MethodGet, server.URL, bytes.NewBufferString("test body"))
	request.Header.Set("Test", "1")
	client := &http.Client{}
	response, _ := client.Do(request)
	responseEntry, err := newResponseEntry(response)

	if responseEntry.StatusCode != 200 {
		t.Error("expect status code to be 200")
	}

	if responseEntry.binaryBody {
		t.Error("expect binaryBody to be false but true")
	}

	if err != nil {
		t.Error("expect no errors")
	}

	if val := responseEntry.Header.Get("Date"); val == "" {
		t.Error("expect value to be not nil but got empty", val)
	}
}
