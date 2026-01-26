package builder

import "github.com/aws/aws-lambda-go/events"

func SQSEvent(body string) events.SQSEvent {
	return events.SQSEvent{
		Records: []events.SQSMessage{
			{
				MessageId:     "example-message-id",
				ReceiptHandle: "example-receipt-handle",
				Body:          body,
			},
		},
	}
}
