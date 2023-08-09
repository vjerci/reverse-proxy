package server_test

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vjerci/reverse-proxy/internal/block"
	"github.com/vjerci/reverse-proxy/internal/log"
	"github.com/vjerci/reverse-proxy/internal/mask"
	"github.com/vjerci/reverse-proxy/internal/proxy"
	"github.com/vjerci/reverse-proxy/internal/server"
)

type InspectorMock struct {
	method func(bytes []byte) ([]byte, error)
}

func (inspector *InspectorMock) Inspect(bytes []byte) ([]byte, error) {
	return inspector.method(bytes)
}

type LoggerMock struct{}

func (logger *LoggerMock) Print(...any) {}

type GuardMock struct {
	method func(req *http.Request) bool
}

func (guard *GuardMock) ShouldBlock(req *http.Request) bool {
	return guard.method(req)
}

type ProxyMock struct {
	method func(req *http.Request, host string, scheme string) (*http.Response, error)
}

func (proxy *ProxyMock) Forward(req *http.Request, host string, scheme string) (*http.Response, error) {
	return proxy.method(req, host, scheme)
}

type BodyErrReaderMock struct{}

func (body *BodyErrReaderMock) Read(p []byte) (n int, err error) {
	return 1, errors.New("mock error")
}

func TestHandleErrors(t *testing.T) {
	const forwardHost = "api.domain.com"
	const forwardScheme = "https"
	const url = "http://localhost:8000"

	var errorHeaders = http.Header{}
	errorHeaders.Set(server.ProxyResponseHeader, server.ProxyResponseHeaderError)

	jsonHeaders := http.Header{
		"Content-Type": []string{"application/json"},
	}

	testCases := []struct {
		testName        string
		expectedStatus  int
		expectedContent []byte

		Inspector             mask.Inspector
		ResponseWriterFactory log.ResponseWriterFactory
		Guard                 block.Guard
		Proxy                 proxy.Proxy

		req  *http.Request
		resp httptest.ResponseRecorder
	}{
		{
			testName:        "body_error",
			expectedStatus:  http.StatusInternalServerError,
			expectedContent: server.ProxyErrorReadingRequestBody,

			ResponseWriterFactory: &log.ResponseWriterFactoryInstance{
				Logger: &LoggerMock{},
			},
			Inspector: nil,
			Guard:     nil,
			Proxy:     nil,

			req:  httptest.NewRequest(http.MethodPost, url, &BodyErrReaderMock{}),
			resp: *httptest.NewRecorder(),
		},
		{
			testName:        "config_guard_block",
			expectedStatus:  http.StatusForbidden,
			expectedContent: server.ProxyErrorBlock,

			ResponseWriterFactory: &log.ResponseWriterFactoryInstance{
				Logger: &LoggerMock{},
			},
			Guard: &GuardMock{
				func(req *http.Request) bool {
					return true
				},
			},
			Inspector: nil,
			Proxy:     nil,

			req:  httptest.NewRequest(http.MethodPost, url, strings.NewReader("")),
			resp: *httptest.NewRecorder(),
		},
		{
			testName:        "proxy_forward_error",
			expectedStatus:  http.StatusInternalServerError,
			expectedContent: server.ProxyErrorForwardingRequest,

			ResponseWriterFactory: &log.ResponseWriterFactoryInstance{
				Logger: &LoggerMock{},
			},
			Guard: &GuardMock{
				func(req *http.Request) bool {
					return false
				},
			},
			Proxy: &ProxyMock{
				method: func(req *http.Request, host string, scheme string) (*http.Response, error) {
					return nil, errors.New("test error")
				},
			},
			Inspector: nil,

			req:  httptest.NewRequest(http.MethodPost, url, strings.NewReader("")),
			resp: *httptest.NewRecorder(),
		},
		{
			testName:        "inspect_error",
			expectedStatus:  http.StatusInternalServerError,
			expectedContent: server.ProxyErrorInspectingRequest,

			ResponseWriterFactory: &log.ResponseWriterFactoryInstance{
				Logger: &LoggerMock{},
			},
			Guard: &GuardMock{
				func(req *http.Request) bool {
					return false
				},
			},
			Proxy: &ProxyMock{
				method: func(req *http.Request, host string, scheme string) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader("{}")),
						Header:     jsonHeaders,
					}, nil
				},
			},
			Inspector: &InspectorMock{
				method: func(bytes []byte) ([]byte, error) {
					return nil, errors.New("test error")
				},
			},

			req:  httptest.NewRequest(http.MethodGet, url, strings.NewReader("[{]}")),
			resp: *httptest.NewRecorder(),
		},
	}

	for _, test := range testCases {
		handler := server.Handle(test.Inspector, test.ResponseWriterFactory, test.Guard, test.Proxy, forwardHost, forwardScheme)
		handler(&test.resp, test.req)

		assert.Equal(t, test.expectedStatus, test.resp.Result().StatusCode, test.testName+" didnt get expected status code")

		respBytes, err := io.ReadAll(test.resp.Body)
		if err != nil {
			t.Fatalf("%s failed to read response body %s", test.testName, err)
		}

		assert.Equal(t, string(test.expectedContent), string(respBytes), test.testName+" didnt get expected body")

		assert.Equal(t, errorHeaders, test.resp.Header(), test.testName+" didnt get expected headers")
	}
}

func TestHandleSucces(t *testing.T) {
	const forwardHost = "api.domain.com"
	const forwardScheme = "https"
	const url = "http://localhost:8000"

	var successHeaders = http.Header{}
	successHeaders.Set(server.ProxyResponseHeader, server.ProxyResponseHeaderSuccess)

	jsonHeaders := http.Header{
		"Content-Type":   []string{"application/json"},
		"Content-Length": []string{"2"},
	}

	responseBodyContent := "{}"

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(responseBodyContent)),
		Header:     jsonHeaders,
	}

	testCase := struct {
		Inspector             mask.Inspector
		ResponseWriterFactory log.ResponseWriterFactory
		Guard                 block.Guard
		Proxy                 proxy.Proxy

		req  *http.Request
		resp httptest.ResponseRecorder
	}{
		ResponseWriterFactory: &log.ResponseWriterFactoryInstance{
			Logger: &LoggerMock{},
		},
		Guard: &GuardMock{
			func(req *http.Request) bool {
				return false
			},
		},
		Proxy: &ProxyMock{
			method: func(req *http.Request, host string, scheme string) (*http.Response, error) {
				return response, nil
			},
		},
		Inspector: &InspectorMock{
			method: func(bytes []byte) ([]byte, error) {
				return bytes, nil
			},
		},

		req:  httptest.NewRequest(http.MethodGet, url, strings.NewReader("{}")),
		resp: *httptest.NewRecorder(),
	}

	handler := server.Handle(testCase.Inspector, testCase.ResponseWriterFactory, testCase.Guard, testCase.Proxy, forwardHost, forwardScheme)
	handler(&testCase.resp, testCase.req)

	assert.Equal(t, response.StatusCode, testCase.resp.Result().StatusCode, "didnt get expected status code")

	respBytes, err := io.ReadAll(testCase.resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body %s", err)
	}

	assert.Equal(t, responseBodyContent, string(respBytes), "didnt get expected body")

	jsonHeaders[server.ProxyResponseHeader] = []string{server.ProxyResponseHeaderSuccess}

	assert.Equal(t, jsonHeaders, testCase.resp.Header(), "headers are invalid")
}
