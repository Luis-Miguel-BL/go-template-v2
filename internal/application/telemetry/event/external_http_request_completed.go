package event

import "time"

type ExternalHTTPRequestCompleted struct {
	Method     string
	URL        string
	StatusCode int
	Duration   time.Duration
}

func (e ExternalHTTPRequestCompleted) Name() string {
	return "ExternalHTTPRequestCompleted"
}

func (e ExternalHTTPRequestCompleted) Attributes() map[string]any {
	success := false
	if e.StatusCode >= 200 && e.StatusCode < 300 {
		success = true
	}

	return map[string]any{
		"method":      e.Method,
		"url":         e.URL,
		"status_code": e.StatusCode,
		"duration_ms": e.Duration.Milliseconds(),
		"success":     success,
	}
}
