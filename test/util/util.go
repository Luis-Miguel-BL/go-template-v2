package util

import (
	"fmt"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/httpclient"
)

type TestUtil struct {
	cfg        *config.Config
	httpClient httpclient.Client
	*aws.SQSClient
}

func NewTestUtil(cfg *config.Config, httpClientFactory httpclient.HTTPClientFactory, sqsClient *aws.SQSClient) *TestUtil {
	baseURL := fmt.Sprintf("http://localhost:%d", cfg.Server.Port)
	return &TestUtil{
		cfg:        cfg,
		httpClient: httpClientFactory.New(httpclient.WithBaseURL(baseURL)),
		SQSClient:  sqsClient,
	}
}

func (t *TestUtil) GetConfig() *config.Config {
	return t.cfg
}
