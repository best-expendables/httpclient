package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout = 3 * time.Second
)

type ResponseParser interface {
	Parse(r io.Reader, result interface{}) error
}

type HeaderSetterFn func(r *http.Request)

type DefaultApiResponseParser struct{}

func (p *DefaultApiResponseParser) Parse(r io.Reader, result interface{}) error {
	container := ApiResponse{
		Data: result,
	}
	return json.NewDecoder(r).Decode(&container)
}

// ApiResponse data structure to extract restful response
type ApiResponse struct {
	Data interface{} `json:"data"`
}

type BaseClient struct {
	baseUrl        string
	timeout        time.Duration
	transport      http.RoundTripper
	httpClient     *http.Client
	responseParser ResponseParser
	headerSetterFn HeaderSetterFn
}

type option func(client *BaseClient)

// WithTimeout timeout option
func WithTimeout(timeout time.Duration) option {
	return func(client *BaseClient) {
		client.timeout = timeout
	}
}

// WithTransport custom transport option
func WithTransport(rt http.RoundTripper) option {
	return func(client *BaseClient) {
		client.transport = rt
	}
}

// WithResponseParser custom response parsing format
func WithResponseParser(p ResponseParser) option {
	return func(client *BaseClient) {
		client.responseParser = p
	}
}

// WithHeaders
func WithHeaderSetterFn(setterFn HeaderSetterFn) option {
	return func(client *BaseClient) {
		client.headerSetterFn = setterFn
	}
}

func NewBaseClient(url string, opts ...option) *BaseClient {
	c := &BaseClient{
		baseUrl:        strings.TrimRight(url, "/"),
		timeout:        defaultTimeout,
		responseParser: &DefaultApiResponseParser{},
	}
	for _, opt := range opts {
		opt(c)
	}
	c.httpClient = &http.Client{
		Timeout:   c.timeout,
		Transport: c.transport,
	}
	return c
}

func (c *BaseClient) DoRequest(ctx context.Context, method, path string, queryParams url.Values, body, result interface{}) error {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	url := c.buildURL(path, queryParams)
	req, err = c.buildRequest(method, url, body)
	if err != nil {
		return err
	}

	resp, err = c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	err = c.checkResponseError(resp)
	if err != nil {
		return err
	}

	if result == nil {
		return nil
	}

	if resp.ContentLength > 0 {
		return c.responseParser.Parse(resp.Body, result)
	}

	// in case content-length is not set correctly but we have data in body
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %s", err.Error())
	}
	if len(content) > 0 {
		return c.responseParser.Parse(bytes.NewReader(content), result)
	}

	// get from header
	if location := resp.Header.Get("Location"); len(location) > 0 {
		if s, ok := result.(*string); ok {
			*s = location
			return nil
		}
		return fmt.Errorf("wrong data type for receive result, expect string pointer but got %T", result)
	}
	return nil
}

func (c *BaseClient) buildRequest(method, url string, body interface{}) (req *http.Request, err error) {
	if body == nil {
		req, err = http.NewRequest(method, url, nil)
	} else {
		var buf *bytes.Buffer
		buf = bytes.NewBuffer([]byte{})
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, buf)
	}
	if err != nil {
		return nil, err
	}
	// add default headers
	req.Header.Add("Content-Type", "application/json")
	if c.headerSetterFn != nil {
		c.headerSetterFn(req)
	}
	return req, nil
}

func (c *BaseClient) buildURL(path string, queryParams url.Values) string {
	buf := bytes.NewBufferString(c.baseUrl)
	buf.WriteByte('/')
	path = strings.TrimLeft(path, "/")
	buf.WriteString(path)
	if len(queryParams) > 0 {
		buf.WriteByte('?')
		buf.WriteString(queryParams.Encode())
	}
	return buf.String()
}

func (c *BaseClient) checkResponseError(res *http.Response) error {
	if res.StatusCode >= 400 {
		err := &Error{
			Code: res.StatusCode,
		}
		if res.ContentLength > 0 {
			errDetail, _ := ioutil.ReadAll(res.Body)
			err.Message = string(errDetail)
		}
		return err
	}
	return nil
}
