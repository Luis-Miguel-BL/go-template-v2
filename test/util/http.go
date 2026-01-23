package util

import (
	"context"
	"net/http"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/http/response"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/httpclient"
)

type RequestParams struct {
	Method      string
	Path        string
	Body        any
	AccessToken string
	Headers     map[string]string
}

func (t *TestUtil) DoRequest(ctx context.Context, params RequestParams) (httpclient.Response, error) {
	options := []httpclient.Option{}
	if params.Body != nil {
		options = append(options, httpclient.WithBody(params.Body))
	}
	for k, v := range params.Headers {
		options = append(options, httpclient.WithHeader(k, v))
	}
	if params.AccessToken != "" {
		options = append(options, httpclient.WithHeader("Authorization", "Bearer "+params.AccessToken))
	}
	options = append(options, httpclient.WithHeader("Content-Type", "application/json"))

	resp, err := t.httpClient.Do(
		ctx,
		params.Method,
		params.Path,
		options...,
	)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *TestUtil) RequestWithAppKey(ctx context.Context, params RequestParams) (httpclient.Response, error) {
	if params.Headers == nil {
		params.Headers = make(map[string]string)
	}
	params.Headers["x-api-key"] = t.cfg.Server.AppKey

	return t.DoRequest(ctx, params)
}

func (t *TestUtil) RequestWithDefaultAccessToken(ctx context.Context, params RequestParams) (httpclient.Response, error) {
	resp, err := t.RequestWithAppKey(ctx, RequestParams{
		Method: http.MethodPost,
		Path:   "/v1/authorization",
	})
	if err != nil {
		return nil, err
	}

	var authResponse response.AuthorizationResponse
	if err := resp.Unmarshal(&authResponse); err != nil {
		return nil, err
	}

	if params.Headers == nil {
		params.Headers = make(map[string]string)
	}
	params.AccessToken = authResponse.AccessToken

	return t.DoRequest(ctx, params)
}
