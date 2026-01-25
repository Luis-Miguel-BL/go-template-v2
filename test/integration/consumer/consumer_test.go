package consumer

import (
	"sync"
	"testing"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/integration"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/integration/mocks"
	_fx "github.com/Luis-Miguel-BL/go-lm-template/internal/fx"
	"github.com/Luis-Miguel-BL/go-lm-template/test/setup"
	"github.com/stretchr/testify/suite"
)

type ConsumerSuite struct {
	setup.BaseTestSuite
	exampleAPIIntegration *mocks.ExampleAPIIntegration
}

func (s *ConsumerSuite) SetupSuite() {
	wg := &sync.WaitGroup{}

	s.exampleAPIIntegration = &mocks.ExampleAPIIntegration{}
	s.SetupMock(s.exampleAPIIntegration, new(integration.ExampleAPIIntegration))
	s.SetupApp(wg,
		_fx.RootModule,
		_fx.ConsumerModule(wg),
		_fx.ApplicationModule(wg),
	)
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(ConsumerSuite))
}
