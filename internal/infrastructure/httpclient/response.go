package httpclient

import (
	"net/http"
	"time"
)

type Response interface {
	Raw() *http.Response
	StatusCode() int
	Body() []byte
	IsSuccess() bool
	Unmarshal(v any) error
	Duration() time.Duration
}
