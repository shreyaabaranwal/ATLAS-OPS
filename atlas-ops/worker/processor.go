package worker

import (
	"context"
	"log"
	"time"


	"atlas-ops/execution"
	"atlas-ops/incident"
	"atlas-ops/queue"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func StartWorker(cfg sdkaws.Config, queueURL string, store *incident.DynamoStore) {

	sqsClient := queue.NewSQSClient(cfg, queueURL)
	scaler := execution.NewEC2Scaler(cfg)

	log.Println("👷 Worker started...")

	for {

		output, err := sqsClient.Client.ReceiveMessage(context.Background(), &sqs.ReceiveMessageInput{
			QueueUrl:            &queueURL,
			MaxNumberOfMessages: 1,
			WaitTimeSeconds:     10,
		})

		if err != nil {
			log.Println("SQS receive error:", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, msg := range output.Messages {

			incidentID := *msg.Body
			log.Println("Processing incident:", incidentID)

			inc, err := store.Get(context.Background(), incidentID)
			if err != nil {
				log.Println("Incident fetch error:", err)
				continue
			}

			targetType := "t3.medium"

			// Execute scaling
			scaler.DryRun(inc.InstanceID, targetType)
			scaler.Execute(inc.InstanceID, targetType)

			// Mark verified
			inc.State = incident.Verified
			store.Create(context.Background(), inc)

			// Delete message from queue
			_, err = sqsClient.Client.DeleteMessage(context.Background(), &sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})

			if err != nil {
				log.Println("Delete message error:", err)
			}
		}
	}
}