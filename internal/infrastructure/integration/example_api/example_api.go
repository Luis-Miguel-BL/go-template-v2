package exampleapi

import (
	"context"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/errors"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/config"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/httpclient"
)

type ExampleAPIIntegration struct {
	httpClient httpclient.Client
}

func NewExampleAPIIntegration(cfg *config.Config, telemetry telemetry.Telemetry, httpClientFactory httpclient.HTTPClientFactory) *ExampleAPIIntegration {

	defaultHeaders := map[string]string{
		"Authorization": "Bearer " + cfg.Integration.ExampleAPI.APIKey,
		"Content-Type":  "application/json",
	}
	httpClient := httpClientFactory.New(
		httpclient.WithBaseURL(cfg.Integration.ExampleAPI.BaseURL),
		httpclient.WithDefaultHeaders(defaultHeaders),
		httpclient.WithMonitoring(telemetry),
	)

	return &ExampleAPIIntegration{
		httpClient: httpClient,
	}
}

func (e *ExampleAPIIntegration) Create(ctx context.Context) (exampleID string, err error) {
	reqBody := CreateExampleRequestDTO{
		// Fill in request body fields as needed
	}
	res, err := e.httpClient.Post(ctx, "/create-endpoint", httpclient.WithBody(reqBody))
	if err != nil {
		return "", errors.UpstreamUnavailableError(0, err)
	}

	if !res.IsSuccess() {
		return "", errors.UpstreamRejectedError(res.StatusCode(), err)
	}

	var resDTO CreateExampleResponseDTO
	err = res.Unmarshal(&resDTO)
	if err != nil {
		return "", errors.UpstreamInvalidResponseError(res.StatusCode(), err)
	}

	return resDTO.ID, nil
}
