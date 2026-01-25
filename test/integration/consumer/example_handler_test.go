package consumer

import (
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/test/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/mock"
)

func (s *ConsumerSuite) TestExampleHandler() {
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
		_, err := s.TestUtil.SQSClient.SendMessage(s.T().Context(), &sqs.SendMessageInput{
			QueueUrl:    &s.TestUtil.GetConfig().Consumer.SQSQueueURL,
			MessageBody: aws.String(msgBody),
		})
		s.NoError(err)
		s.NoError(util.AwaitDone(done, 5*time.Second))
		s.exampleAPIIntegration.AssertExpectations(s.T())
	})
}
