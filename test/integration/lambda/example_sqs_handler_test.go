package lambda

import (
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/lambda/handler"
	"github.com/Luis-Miguel-BL/go-lm-template/test/builder"
	"github.com/Luis-Miguel-BL/go-lm-template/test/util"
	"github.com/stretchr/testify/mock"
)

func (s *LambdaSuite) TestExampleSQSHandler() {
	s.Run("Should process message", func() {
		done := make(chan struct{})

		s.exampleAPIIntegration.
			On("Create", mock.Anything).
			Run(util.DoneChan(done)).
			Return("", nil).Once()

		msgBody := `{
			"id": "test",
			"name": "test",
			"email": "fake@email.com"
		}`

		sqsEvent := builder.SQSEvent(msgBody)
		_, err := s.lambdaRunner.RunLocal(handler.ExampleSQSHandlerLambdaName, sqsEvent)
		s.NoError(err)

		s.NoError(util.AwaitDone(done, 5*time.Second))
		s.exampleAPIIntegration.AssertExpectations(s.T())
	})
}
