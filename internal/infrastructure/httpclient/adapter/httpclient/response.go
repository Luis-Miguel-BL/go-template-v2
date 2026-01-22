package httpclient

import (
	"encoding/json"
	"net/http"
	"time"
)

type Response struct {
	raw        *http.Response
	body       []byte
	statusCode int
	duration   time.Duration
}

func (r *Response) Raw() *http.Response {
	return r.raw
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) IsSuccess() bool {
	return r.statusCode >= 200 && r.statusCode < 300

}

func (r *Response) Unmarshal(v any) error {
	err := json.Unmarshal(r.body, v)
	if err != nil {
		return err
	}
	return nil

}

func (r *Response) Duration() time.Duration {
	return r.duration
}
