package util

import (
	"fmt"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/httpclient"
)

type TestUtil struct {
	cfg        *config.Config
	httpClient httpclient.Client
}

func NewTestUtil(cfg *config.Config, httpClientFactory httpclient.HTTPClientFactory) *TestUtil {
	baseURL := fmt.Sprintf("http://localhost:%d", cfg.Server.Port)
	return &TestUtil{
		cfg:        cfg,
		httpClient: httpClientFactory.New(httpclient.WithBaseURL(baseURL)),
	}
}

func (t *TestUtil) GetConfig() *config.Config {
	return t.cfg
}
