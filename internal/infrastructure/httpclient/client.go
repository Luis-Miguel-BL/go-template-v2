package httpclient

import (
	"context"
)

type Client interface {
	Get(ctx context.Context, endpoint string, opts ...Option) (Response, error)
	Post(ctx context.Context, endpoint string, opts ...Option) (Response, error)
	Put(ctx context.Context, endpoint string, opts ...Option) (Response, error)
	Patch(ctx context.Context, endpoint string, opts ...Option) (Response, error)
	Delete(ctx context.Context, endpoint string, opts ...Option) (Response, error)

	Do(ctx context.Context, method, endpoint string, opts ...Option) (Response, error)
}
