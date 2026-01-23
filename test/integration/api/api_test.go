package api

import (
	"sync"
	"testing"

	_fx "github.com/Luis-Miguel-BL/go-lm-template/internal/fx"
	"github.com/Luis-Miguel-BL/go-lm-template/test/setup"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	setup.BaseTestSuite
}

func (s *APISuite) SetupSuite() {
	wg := &sync.WaitGroup{}
	s.SetupApp(wg,
		_fx.RootModule,
		_fx.HttpModule(wg),
		_fx.ApplicationModule(wg),
	)
}

func TestAPISuite(t *testing.T) {
	suite.Run(t, new(APISuite))
}
