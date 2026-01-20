package metric

import "github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"

type LeadCounter struct {
}

func (c LeadCounter) Name() string {
	return "lead_counter"
}

func (c LeadCounter) Attributes() map[string]any {
	return map[string]any{}
}

func (c LeadCounter) Type() observability.MetricType {
	return observability.MetricTypeCounter
}
func (c LeadCounter) Value() int64 {
	return 1
}
