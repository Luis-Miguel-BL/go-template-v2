package httpclient

type HTTPClientFactory interface {
	New(opts ...ClientOptions) Client
}
