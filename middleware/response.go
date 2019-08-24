package middleware

import (
	"bitbucket.org/snapmartinc/httpclient/net"
	"bytes"
	"io/ioutil"
	"net/http"
)

// Request entry
type responseEntry struct {
	Header     http.Header
	Body       string
	StatusCode int
	binaryBody bool
}

func newResponseEntry(response *http.Response) (responseEntry, error) {
	header := response.Header

	responseEntry := responseEntry{
		Header:     header,
		StatusCode: response.StatusCode,
	}

	if response.Body == nil {
		return responseEntry, nil
	}

	body, err := ioutil.ReadAll(response.Body)
	response.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err != nil {
		return responseEntry, err
	}

	if net.HasBinaryContent(response.Header, body) {
		responseEntry.binaryBody = true

		return responseEntry, nil
	}

	responseEntry.Body = string(body)

	return responseEntry, nil
}
