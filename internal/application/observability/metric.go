package observability

type Metric interface {
	Name() string
	Attributes() map[string]any
	Type() MetricType
	Value() int64
}

type MetricType int

const (
	MetricTypeCounter MetricType = iota
	MetricTypeGauge
	MetricTypeHistogram
)
