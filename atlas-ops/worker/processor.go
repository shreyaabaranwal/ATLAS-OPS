package worker

import (
	"context"
	"log"
	"time"

	"atlas-ops/aws"
	"atlas-ops/execution"
	"atlas-ops/incident"
	"atlas-ops/queue"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

const MaxExecutionAttempts = 3

func StartWorker(cfg sdkaws.Config, queueURL string, store *incident.DynamoStore, audit *incident.AuditStore) {

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

			ctx := context.Background()
			incidentID := *msg.Body

			log.Println("Processing incident:", incidentID)

			// 🔐 Atomic lock
			err := store.TransitionState(ctx, incidentID, incident.Approved, incident.Executing)
			if err != nil {
				log.Println("Another worker already processing or invalid state. Skipping.")
				continue
			}

			// Increment attempts
			_ = store.IncrementAttempts(ctx, incidentID)

			inc, err := store.Get(ctx, incidentID)
			if err != nil {
				log.Println("Fetch error:", err)
				continue
			}

			// 🔥 CIRCUIT BREAKER
			if inc.ExecutionAttempts >= MaxExecutionAttempts {

				log.Println("Circuit breaker triggered")

				_ = store.TransitionState(ctx, inc.ID, incident.Executing, incident.FailedPermanent)

				_ = audit.Save(ctx, &incident.AuditLog{
					IncidentID: inc.ID,
					Action:     "CIRCUIT_BREAKER",
					Actor:      "SYSTEM",
					Result:     "FAILED_PERMANENT",
				})

				// DELETE message because permanent failure
				_, _ = sqsClient.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      &queueURL,
					ReceiptHandle: msg.ReceiptHandle,
				})

				continue
			}

			targetType := "t3.medium" // change to invalid for DLQ test

			err = scaler.Execute(ctx, inc.InstanceID, targetType)
			if err != nil {

				log.Println("Execution failed. Will retry via SQS.")

				_ = store.TransitionState(ctx, inc.ID, incident.Executing, incident.Approved)

				// ❌ DO NOT DELETE MESSAGE
				// Let SQS retry
				continue
			}

			_ = store.TransitionState(ctx, inc.ID, incident.Executing, incident.Executed)

			time.Sleep(30 * time.Second)

			newCPU, err := aws.GetCPUUtilization(cfg, inc.InstanceID)
			if err != nil {
				continue
			}

			if newCPU < inc.CPU {

				now := time.Now()
				inc.VerifiedAt = &now

				_ = store.TransitionState(ctx, inc.ID, incident.Executed, incident.Verified)

				_ = audit.Save(ctx, &incident.AuditLog{
					IncidentID: inc.ID,
					Action:     "VERIFICATION",
					Actor:      "SYSTEM",
					Result:     "VERIFIED",
					CPUBefore:  inc.CPU,
					CPUAfter:   newCPU,
				})

			} else {

				_ = scaler.Rollback(ctx, inc.InstanceID, inc.InstanceType)

				_ = store.TransitionState(ctx, inc.ID, incident.Executed, incident.RolledBack)
			}

			// ✅ Delete only on success path
			_, _ = sqsClient.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})

			log.Println("✅ Incident lifecycle completed safely.")
		}
	}
}