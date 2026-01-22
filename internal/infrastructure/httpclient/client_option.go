package httpclient

import (
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/telemetry"
)

type NewClientOptions struct {
	BaseURL        string
	DefaultHeaders map[string]string
	Timeout        time.Duration
	MonitorEnabled bool
	Telemetry      telemetry.Telemetry
}

type ClientOptions func(*NewClientOptions)

func WithBaseURL(baseURL string) ClientOptions {
	return func(o *NewClientOptions) {
		o.BaseURL = baseURL
	}
}

func WithDefaultHeaders(headers map[string]string) ClientOptions {
	return func(o *NewClientOptions) {
		o.DefaultHeaders = headers
	}
}

func WithDefaultTimeout(timeout time.Duration) ClientOptions {
	return func(o *NewClientOptions) {
		o.Timeout = timeout
	}
}

func WithMonitoring(telemetry telemetry.Telemetry) ClientOptions {
	return func(o *NewClientOptions) {
		o.MonitorEnabled = true
		o.Telemetry = telemetry
	}
}
