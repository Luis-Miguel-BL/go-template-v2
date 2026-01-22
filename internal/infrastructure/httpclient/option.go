package httpclient

import (
	"encoding/json"
	"time"
)

type Option func(*RequestOptions)

type RequestOptions struct {
	Headers map[string]string
	Query   map[string]string
	Body    []byte
	Timeout time.Duration
}

func WithHeader(key, value string) Option {
	return func(o *RequestOptions) {
		o.Headers[key] = value
	}
}

func WithQuery(key, value string) Option {
	return func(o *RequestOptions) {
		o.Query[key] = value
	}
}

func WithBody(body any) Option {
	jsonData, _ := json.Marshal(body)
	return func(o *RequestOptions) {
		o.Body = jsonData
	}
}

func WithBodyBytes(body []byte) Option {
	return func(o *RequestOptions) {
		o.Body = body
	}
}

func WithTimeout(d time.Duration) Option {
	return func(o *RequestOptions) {
		o.Timeout = d
	}
}
