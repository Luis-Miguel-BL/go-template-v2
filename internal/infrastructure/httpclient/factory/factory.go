package factory

import (
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/httpclient"
	_httpclient "github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/httpclient/adapter/httpclient"
)

type httpClientFactory struct {
}

func NewHTTPClientFactory() httpclient.HTTPClientFactory {
	return &httpClientFactory{}
}

func (f *httpClientFactory) New(opts ...httpclient.ClientOptions) httpclient.Client {
	return _httpclient.NewClient(opts...)
}
