package errors

import "fmt"

type ErrorCode string

const (
	ErrCodeUpstreamUnavailable ErrorCode = "UpstreamUnavailable"
	ErrCodeUpstreamInvalid     ErrorCode = "UpstreamInvalidResponse"
	ErrCodeUpstreamRejected    ErrorCode = "UpstreamRejected"
)

type IntegrationError struct {
	Code       ErrorCode
	Err        error
	StatusCode int
}

func (e IntegrationError) Error() string {
	return fmt.Sprintf("integration error [%s]: %s", e.Code, e.Err.Error())
}

func (e IntegrationError) Unwrap() error {
	return e.Err
}

func NewIntegrationError(code ErrorCode, statusCode int, err error) IntegrationError {
	return IntegrationError{
		Code:       code,
		StatusCode: statusCode,
		Err:        err,
	}
}

func UpstreamUnavailableError(statusCode int, err error) IntegrationError {
	return NewIntegrationError(ErrCodeUpstreamUnavailable, statusCode, err)
}

func UpstreamInvalidResponseError(statusCode int, err error) IntegrationError {
	return NewIntegrationError(ErrCodeUpstreamInvalid, statusCode, err)
}

func UpstreamRejectedError(statusCode int, err error) IntegrationError {
	return NewIntegrationError(ErrCodeUpstreamRejected, statusCode, err)
}
