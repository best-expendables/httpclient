package middleware

import (
	"github.com/best-expendables/httpclient/net"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Request entry
type requestEntry struct {
	Header     http.Header
	Method     string
	Body       string
	binaryBody bool
}

func newRequestEntry(request *http.Request) (requestEntry, error) {
	header := request.Header

	requestEntry := requestEntry{
		Method: request.Method,
		Header: header,
	}

	if request.Body == nil {
		return requestEntry, nil
	}

	body, err := ioutil.ReadAll(request.Body)
	defer recoverRequestBody(request, body)

	if err != nil {
		return requestEntry, err
	}

	if net.HasBinaryContent(request.Header, body) {
		requestEntry.binaryBody = true

		return requestEntry, nil
	}

	if contentType := header.Get("Content-Type"); contentType == "application/x-www-form-urlencoded" {
		requestEntry.Body = parseFormToJson(request, body)

		return requestEntry, nil
	}

	requestEntry.Body = string(body)

	return requestEntry, nil
}

func parseFormToJson(request *http.Request, body []byte) string {
	recoverRequestBody(request, body)
	request.ParseForm()

	fields, _ := json.Marshal(request.Form)

	return string(fields)
}
