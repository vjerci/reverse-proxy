package server

import (
	"bytes"
	"io"
	"net/http"
	"strconv"

	"github.com/vjerci/reverse-proxy/internal/block"
	"github.com/vjerci/reverse-proxy/internal/log"
	"github.com/vjerci/reverse-proxy/internal/mask"
	"github.com/vjerci/reverse-proxy/internal/proxy"
)

var ProxyErrorBlock = []byte("proxy config blocks this request")
var ProxyErrorForwardingRequest = []byte("proxy failed to forward request and get response")
var ProxyErrorReadingResponseBody = []byte("proxy failed to read forwarded response body")
var ProxyErrorReadingRequestBody = []byte("proxy failed to read request body")
var ProxyErrorInspectingRequest = []byte("proxy failed to inspect forwarded response body")

const ProxyResponseHeader = "X-Proxy-Error"
const ProxyResponseHeaderError = "true"
const ProxyResponseHeaderSuccess = "false"

func Handle(inspector mask.Inspector, responseWriterFactory log.ResponseWriterFactory, guard block.Guard, proxy proxy.Proxy, forwardHost string, forwardScheme string) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()

		reqBody, err := io.ReadAll(req.Body)
		if err != nil {
			respWithLog := responseWriterFactory.New(req, nil, w)
			respWithLog.Write(http.StatusInternalServerError, map[string][]string{
				ProxyResponseHeader: {ProxyResponseHeaderError},
			}, ProxyErrorReadingRequestBody)
			return
		}

		respWithLog := responseWriterFactory.New(req, nil, w)
		req.Body = io.NopCloser(bytes.NewBuffer(reqBody))

		if guard.ShouldBlock(req) {
			respWithLog.Write(http.StatusForbidden, map[string][]string{
				ProxyResponseHeader: {ProxyResponseHeaderError},
			}, ProxyErrorBlock)
			return
		}

		proxyResp, err := proxy.Forward(req, forwardHost, forwardScheme)
		if err != nil {
			respWithLog.Write(http.StatusInternalServerError, map[string][]string{
				ProxyResponseHeader: {ProxyResponseHeaderError},
			}, ProxyErrorForwardingRequest)
			return
		}

		respBytes, err := io.ReadAll(proxyResp.Body)
		if err != nil {
			respWithLog.Write(http.StatusInternalServerError, map[string][]string{
				ProxyResponseHeader: {ProxyResponseHeaderError},
			}, ProxyErrorReadingResponseBody)
			return
		}

		tamperedRequest := false

		if req.Method == http.MethodGet && proxyResp.Header.Get("Content-Type") == "application/json" {
			maskedJson, err := inspector.Inspect(respBytes)
			if err != nil {
				respWithLog.Write(http.StatusInternalServerError, map[string][]string{
					ProxyResponseHeader: {ProxyResponseHeaderError},
				}, ProxyErrorInspectingRequest)
				return
			}

			respBytes = maskedJson
			tamperedRequest = true

		}

		headers := make(map[string][]string)

		for key, value := range proxyResp.Header {
			headers[key] = value
		}

		if tamperedRequest {
			headers["Content-Length"] = []string{strconv.Itoa(len(respBytes))}
		}

		headers[ProxyResponseHeader] = []string{ProxyResponseHeaderSuccess}

		respWithLog.Write(proxyResp.StatusCode, headers, respBytes)
	}
}
