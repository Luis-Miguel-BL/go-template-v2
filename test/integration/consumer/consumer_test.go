package consumer

import (
	"sync"
	"testing"

	_fx "github.com/Luis-Miguel-BL/go-lm-template/internal/fx"
	"github.com/Luis-Miguel-BL/go-lm-template/test/setup"
	"github.com/stretchr/testify/suite"
)

type ConsumerSuite struct {
	setup.BaseTestSuite
}

func (s *ConsumerSuite) SetupSuite() {
	wg := &sync.WaitGroup{}
	s.SetupApp(wg,
		_fx.RootModule,
		_fx.ConsumerModule(wg),
		_fx.ApplicationModule(wg),
	)
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(ConsumerSuite))
}
