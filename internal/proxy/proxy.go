package proxy

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

var ErrFailedToForward = errors.New("failed to forward request")
var ErrFailedToBuildURL = errors.New("failed to build url")

type Proxy interface {
	Forward(req *http.Request, host string, scheme string) (*http.Response, error)
}

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type ProxyInstance struct {
	http HTTPClient
}

func NewProxy(client HTTPClient) Proxy {
	return &ProxyInstance{
		http: client,
	}
}

func (proxy *ProxyInstance) Forward(req *http.Request, host string, scheme string) (*http.Response, error) {
	req.URL.Host = host
	req.URL.Scheme = scheme
	u, err := url.Parse(req.URL.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToBuildURL, err)
	}
	req.URL = u

	req.RequestURI = ""

	resp, err := proxy.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToForward, err)
	}

	return resp, nil
}
