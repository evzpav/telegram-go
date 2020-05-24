package httpclient

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Mock struct {
	RequestURL     string
	RequestHeader  http.Header
	RequestBody    string
	RequestMethod  string
	ResponseHeader http.Header
	ResponseBody   string
	ResponseStatus int
	Error          error
}

func (m *Mock) Do(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		m.RequestBody = string(body)
	}

	m.RequestMethod = req.Method
	m.RequestURL = req.URL.String()
	m.RequestHeader = req.Header

	if m.Error != nil {
		return nil, m.Error
	}

	response := ioutil.NopCloser(bytes.NewReader([]byte(m.ResponseBody)))

	return &http.Response{
		StatusCode: m.ResponseStatus,
		Body:       response,
		Header:     m.ResponseHeader,
	}, nil
}

func (m *Mock) Status(statusCode int) *Mock {
	m.ResponseStatus = statusCode
	return m
}

func (m *Mock) Body(body string) *Mock {
	m.ResponseBody = body
	return m
}

func (m *Mock) Err(err error) *Mock {
	m.Error = err
	return m
}

type httpClientMultMock struct {
	mocks map[string]map[string]*Mock
}

func NewHTTPMultMock() *httpClientMultMock {
	return &httpClientMultMock{make(map[string]map[string]*Mock)}
}

func (m *httpClientMultMock) createMock(method, URL string) *Mock {
	methodMocks, ok := m.mocks[method]

	if !ok {
		methodMocks = make(map[string]*Mock)
		m.mocks[method] = methodMocks
	}
	mock := &Mock{}
	methodMocks[URL] = mock
	return mock
}

func (m *httpClientMultMock) Get(URL string) *Mock {
	return m.createMock("GET", URL)
}

func (m *httpClientMultMock) Put(URL string) *Mock {
	return m.createMock("PUT", URL)
}

func (m *httpClientMultMock) Post(URL string) *Mock {
	return m.createMock("POST", URL)
}

func (m *httpClientMultMock) Do(req *http.Request) (*http.Response, error) {
	method := req.Method
	URL := req.URL.String()

	methodMocks, ok := m.mocks[method]
	if !ok {
		return nil, fmt.Errorf("No mock for [%s][%s]", method, URL)
	}
	pathMock, ok := methodMocks[URL]
	if !ok {
		return nil, fmt.Errorf("No mock for method [%s][%s]", method, URL)
	}

	return pathMock.Do(req)
}
