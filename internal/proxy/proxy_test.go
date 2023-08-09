package proxy_test

import (
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/vjerci/reverse-proxy/internal/proxy"
)

type HTTPClientMock struct {
	Req  *http.Request
	Err  error
	Resp *http.Response
}

func (client *HTTPClientMock) Do(req *http.Request) (*http.Response, error) {
	client.Req = req
	return client.Resp, client.Err
}

func TestHttpProxyError(t *testing.T) {
	proxyInstance := proxy.NewProxy(&HTTPClientMock{
		Err: errors.New("dummy error"),
	})

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/api", strings.NewReader(""))
	if err != nil {
		t.Fatalf("failed to create request %s", err)
	}

	_, err = proxyInstance.Forward(req, "api.domain.com", "https")

	if !errors.Is(err, proxy.ErrFailedToForward) {
		t.Fatalf("expected to get errFailedToForward err got '%s' instead", err)
	}
}

func TestHttpProxySuccess(t *testing.T) {
	clientMock := &HTTPClientMock{
		Resp: &http.Response{},
	}
	proxyInstance := proxy.NewProxy(clientMock)

	forwardHost := "api.domain.com"
	forwardScheme := "https"

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8000/api", strings.NewReader(""))
	if err != nil {
		t.Fatalf("failed to create request %s", err)
	}

	resp, err := proxyInstance.Forward(req, forwardHost, forwardScheme)

	if err != nil {
		t.Fatalf("expected success got err instead %s", err)
	}

	if resp == nil {
		t.Fatal("expected resp got nil instead")
	}

	if clientMock.Req.URL.Host != forwardHost {
		t.Fatalf("expected to replace host on forwarding request  got %s instead", clientMock.Req.URL.Host)
	}
}
