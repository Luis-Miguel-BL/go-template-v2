package lambda

import (
	"sync"
	"testing"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/integration"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/integration/mocks"
	_fx "github.com/Luis-Miguel-BL/go-lm-template/internal/fx"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/lambda"
	"github.com/Luis-Miguel-BL/go-lm-template/test/setup"
	"github.com/stretchr/testify/suite"
	"go.uber.org/fx"
)

type LambdaSuite struct {
	setup.BaseTestSuite
	exampleAPIIntegration *mocks.ExampleAPIIntegration
	lambdaRunner          *lambda.Runner
}

func (s *LambdaSuite) SetupSuite() {
	wg := &sync.WaitGroup{}

	s.exampleAPIIntegration = &mocks.ExampleAPIIntegration{}
	s.SetupMock(s.exampleAPIIntegration, new(integration.ExampleAPIIntegration))
	s.SetupApp(wg,
		_fx.RootModule,
		_fx.LambdaModule,
		_fx.ApplicationModule(wg),

		fx.Invoke(func(runner *lambda.Runner) {
			s.lambdaRunner = runner
		}),
	)
}

func TestLambdaSuite(t *testing.T) {
	suite.Run(t, new(LambdaSuite))
}
