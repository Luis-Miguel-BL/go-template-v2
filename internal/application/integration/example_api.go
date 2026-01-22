package integration

import "context"

type ExampleAPIIntegration interface {
	Create(ctx context.Context) (exampleID string, err error)
}
