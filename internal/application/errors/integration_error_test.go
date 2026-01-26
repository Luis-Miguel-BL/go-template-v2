package errors

import (
	"errors"
	"testing"
)

func TestIntegrationError_Error(t *testing.T) {
	baseErr := errors.New("connection refused")

	err := IntegrationError{
		Code:       ErrCodeUpstreamUnavailable,
		StatusCode: 503,
		Err:        baseErr,
	}

	expected := "integration error [UpstreamUnavailable]: connection refused"

	if err.Error() != expected {
		t.Errorf("expected %q, got %q", expected, err.Error())
	}
}

func TestIntegrationError_Unwrap(t *testing.T) {
	baseErr := errors.New("timeout")

	err := IntegrationError{
		Code:       ErrCodeUpstreamInvalid,
		StatusCode: 502,
		Err:        baseErr,
	}

	if !errors.Is(err, baseErr) {
		t.Errorf("expected unwrap to return base error")
	}
}

func TestNewIntegrationError(t *testing.T) {
	baseErr := errors.New("invalid payload")

	err := NewIntegrationError(
		ErrCodeUpstreamRejected,
		400,
		baseErr,
	)

	if err.Code != ErrCodeUpstreamRejected {
		t.Errorf("expected code %v, got %v", ErrCodeUpstreamRejected, err.Code)
	}

	if err.StatusCode != 400 {
		t.Errorf("expected status code 400, got %d", err.StatusCode)
	}

	if err.Err != baseErr {
		t.Errorf("expected base error %v, got %v", baseErr, err.Err)
	}
}

func TestUpstreamUnavailableError(t *testing.T) {
	baseErr := errors.New("service down")

	err := UpstreamUnavailableError(503, baseErr)

	if err.Code != ErrCodeUpstreamUnavailable {
		t.Errorf("expected code %v, got %v", ErrCodeUpstreamUnavailable, err.Code)
	}

	if err.StatusCode != 503 {
		t.Errorf("expected status code 503, got %d", err.StatusCode)
	}

	if err.Err != baseErr {
		t.Errorf("expected base error %v, got %v", baseErr, err.Err)
	}
}

func TestUpstreamInvalidResponseError(t *testing.T) {
	baseErr := errors.New("invalid json")

	err := UpstreamInvalidResponseError(502, baseErr)

	if err.Code != ErrCodeUpstreamInvalid {
		t.Errorf("expected code %v, got %v", ErrCodeUpstreamInvalid, err.Code)
	}
}

func TestUpstreamRejectedError(t *testing.T) {
	baseErr := errors.New("unauthorized")

	err := UpstreamRejectedError(401, baseErr)

	if err.Code != ErrCodeUpstreamRejected {
		t.Errorf("expected code %v, got %v", ErrCodeUpstreamRejected, err.Code)
	}
}
