package httpclient

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry/event"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/httpclient"
)

const defaultTimeout = 10 * time.Second

type client struct {
	baseURL        string
	http           *http.Client
	headers        map[string]string
	monitorEnabled bool
	telemetry      telemetry.Telemetry
}

func NewClient(opts ...httpclient.ClientOptions) *client {
	options := &httpclient.NewClientOptions{
		Timeout:        defaultTimeout,
		DefaultHeaders: make(map[string]string),
	}

	for _, opt := range opts {
		opt(options)
	}
	httpClient := &http.Client{
		Timeout: options.Timeout,
	}

	if options.MonitorEnabled && options.Telemetry != nil {
		httpClient.Transport = options.Telemetry.NewHttpTransport()
	}

	return &client{
		baseURL:        options.BaseURL,
		http:           httpClient,
		headers:        options.DefaultHeaders,
		monitorEnabled: options.MonitorEnabled,
		telemetry:      options.Telemetry,
	}
}

func (c *client) Get(ctx context.Context, endpoint string, opts ...httpclient.Option) (httpclient.Response, error) {
	return c.Do(ctx, http.MethodGet, endpoint, opts...)
}

func (c *client) Post(ctx context.Context, endpoint string, opts ...httpclient.Option) (httpclient.Response, error) {
	return c.Do(ctx, http.MethodPost, endpoint, opts...)
}

func (c *client) Put(ctx context.Context, endpoint string, opts ...httpclient.Option) (httpclient.Response, error) {
	return c.Do(ctx, http.MethodPut, endpoint, opts...)
}

func (c *client) Patch(ctx context.Context, endpoint string, opts ...httpclient.Option) (httpclient.Response, error) {
	return c.Do(ctx, http.MethodPatch, endpoint, opts...)
}

func (c *client) Delete(ctx context.Context, endpoint string, opts ...httpclient.Option) (httpclient.Response, error) {
	return c.Do(ctx, http.MethodDelete, endpoint, opts...)
}

func (c *client) Do(ctx context.Context, method, endpoint string, opts ...httpclient.Option) (httpclient.Response, error) {
	var span telemetry.Span
	if c.monitorEnabled {
		ctx, span = c.telemetry.StartSpan(ctx, "http.client.request")
		defer span.End()
	}

	ro := &httpclient.RequestOptions{
		Headers: map[string]string{},
		Query:   map[string]string{},
		Timeout: c.http.Timeout,
	}

	for _, opt := range opts {
		opt(ro)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.baseURL+endpoint,
		bytes.NewBuffer(ro.Body),
	)
	if err != nil {
		return nil, err
	}

	// default headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// request headers
	for k, v := range ro.Headers {
		req.Header.Set(k, v)
	}

	// query params
	q := req.URL.Query()
	for k, v := range ro.Query {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	start := time.Now()

	ctx, cancel := context.WithTimeout(ctx, ro.Timeout)
	defer cancel()

	resp, err := c.http.Do(req.WithContext(ctx))
	if err != nil {
		c.telemetry.RecordError(ctx, err)
		return nil, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)

	if c.monitorEnabled {
		fullURL := req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
		c.telemetry.AddAttributes(ctx, map[string]any{
			"http.method":      method,
			"http.url":         fullURL,
			"http.status_code": resp.StatusCode,
			"http.duration_ms": duration.Milliseconds(),
		})

		c.telemetry.AddEvent(ctx, event.ExternalHTTPRequestCompleted{
			Method:     method,
			URL:        fullURL,
			StatusCode: resp.StatusCode,
			Duration:   duration,
		})
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.telemetry.RecordError(ctx, err)
		return nil, err
	}

	return &Response{
		raw:        resp,
		body:       body,
		statusCode: resp.StatusCode,
		duration:   duration,
	}, nil
}
