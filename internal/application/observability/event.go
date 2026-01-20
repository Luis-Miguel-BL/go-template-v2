package observability

type Event interface {
	Name() string
	Attributes() map[string]any
}
