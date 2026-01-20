package sqs

import (
	"context"
	"sync"
	"time"

	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/logger"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/application/observability"
	"github.com/Luis-Miguel-BL/go-lm-template/internal/infrastructure/aws"
	_aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Consumer struct {
	cfg     ConsumerConfig
	logger  logger.Logger
	client  *aws.SQSClient
	handler Handler
	ctx     context.Context
	cancel  context.CancelFunc
}

type ConsumerConfig struct {
	QueueURL            string
	DLQURL              string
	WaitTimeSeconds     int32
	MaxNumberOfMessages int32
	VisibilityTimeout   int32
	PollIntervalSeconds int
}

func NewConsumer(
	cfg ConsumerConfig,
	client *aws.SQSClient,
	handler Handler,
	logger logger.Logger,
) *Consumer {
	if cfg.MaxNumberOfMessages == 0 {
		cfg.MaxNumberOfMessages = 10
	}
	if cfg.WaitTimeSeconds == 0 {
		cfg.WaitTimeSeconds = 20
	}
	if cfg.VisibilityTimeout == 0 {
		cfg.VisibilityTimeout = 30
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Consumer{
		cfg:     cfg,
		client:  client,
		handler: handler,
		ctx:     ctx,
		cancel:  cancel,
		logger:  logger.WithFields(map[string]any{"queue_url": cfg.QueueURL}),
	}
}

func (c *Consumer) Start(wg *sync.WaitGroup) {
	c.logger.Info("starting consumer" + c.cfg.QueueURL)
	wg.Add(1)
	c.run(wg)

}

func (c *Consumer) Stop() {
	c.logger.Info("stopping consumer " + c.cfg.QueueURL)
	c.cancel()
}

func (c *Consumer) run(wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-c.ctx.Done():
			return
		default:
			time.Sleep(time.Second)
			output, err := c.client.ReceiveMessage(c.ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            _aws.String(c.cfg.QueueURL),
				MaxNumberOfMessages: c.cfg.MaxNumberOfMessages,
				WaitTimeSeconds:     c.cfg.WaitTimeSeconds,
				VisibilityTimeout:   c.cfg.VisibilityTimeout,
			})
			if err != nil {
				c.logger.Error("receive error", err)
				continue
			}

			for _, msg := range output.Messages {
				result := c.handle(c.ctx, msg)
				err := c.handleResult(c.ctx, msg, result)
				if err != nil {
					c.logger.Error("handle result error", result, err)
					continue
				}
			}
		}
	}
}

func (c *Consumer) handle(ctx context.Context, msg types.Message) HandleResult {
	ctx, span := observability.GetObservability().StartSpan(c.ctx, "sqs.consumer.handle_message")
	defer span.End()

	result, err := c.handler.Handle(ctx, msg)
	if err != nil {
		c.logger.Error("handler error",
			"messageId", msg.MessageId,
			"error", err,
		)
		span.RecordError(err)
	}
	span.SetAttributes(map[string]any{"handle_result": result})

	return result
}

func (c *Consumer) handleResult(ctx context.Context, msg types.Message, result HandleResult) error {
	switch result {
	case HandleSuccess:
		_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      _aws.String(c.cfg.QueueURL),
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			return err
		}

	case HandleRetry:
		c.logger.Warn("message scheduled for retry", "messageId", msg.MessageId)

	case HandleDLQ:
		if c.cfg.DLQURL == "" {
			c.logger.Error("dlq url not configured, cannot send message to dlq", "messageId", msg.MessageId)
			return nil
		}

		_, err := c.client.SendMessage(ctx, &sqs.SendMessageInput{
			QueueUrl:    _aws.String(c.cfg.DLQURL),
			MessageBody: msg.Body,
		})
		if err != nil {
			c.logger.Error("failed to send to dlq", err)
			return err
		}
		_, err = c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
			QueueUrl:      _aws.String(c.cfg.QueueURL),
			ReceiptHandle: msg.ReceiptHandle,
		})
		if err != nil {
			c.logger.Error("delete error", err)
			return err
		}
	}
	return nil
}
