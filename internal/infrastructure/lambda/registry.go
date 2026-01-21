package lambda

import "fmt"

type LambdaHandler interface {
	LambdaName() string
}

type Registry struct {
	handlersMap map[string]LambdaHandler
}

func NewRegistry(handlers []LambdaHandler) *Registry {
	handlersMap := make(map[string]LambdaHandler, len(handlers))

	for _, handler := range handlers {
		handlersMap[handler.LambdaName()] = handler
	}

	return &Registry{
		handlersMap: handlersMap,
	}
}

func (r *Registry) Get(name string) (LambdaHandler, error) {
	h, ok := r.handlersMap[name]
	if !ok {
		return nil, fmt.Errorf("lambda handler '%s' not registered", name)
	}
	return h, nil
}
