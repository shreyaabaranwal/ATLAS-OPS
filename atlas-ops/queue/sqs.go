package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSClient struct {
	Client   *sqs.Client
	QueueURL string
}

func NewSQSClient(cfg aws.Config, queueURL string) *SQSClient {
	return &SQSClient{
		Client:   sqs.NewFromConfig(cfg),
		QueueURL: queueURL,
	}
}

func (s *SQSClient) SendMessage(ctx context.Context, incidentID string) error {

	_, err := s.Client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(incidentID),
	})

	return err
}