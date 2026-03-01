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

			ctx := context.Background()

			incidentID := *msg.Body
			log.Println("Processing incident:", incidentID)

			inc, err := store.Get(ctx, incidentID)
			if err != nil {
				log.Println("Incident fetch error:", err)
				continue
			}

			targetType := "t3.medium"

			// ---------------- 1️⃣ DRY RUN ----------------
			err = scaler.DryRun(ctx, inc.InstanceID, targetType)
			if err != nil {
				log.Println("DryRun failed:", err)

				inc.State = incident.RolledBack
				store.Update(ctx, inc)
				continue
			}

			inc.State = incident.Simulated
			store.Update(ctx, inc)

			// ---------------- 2️⃣ EXECUTE ----------------
			startTime := time.Now()

			err = scaler.Execute(ctx, inc.InstanceID, targetType)
			if err != nil {
				log.Println("Execution failed:", err)

				inc.State = incident.RolledBack
				store.Update(ctx, inc)
				continue
			}

			inc.State = incident.Executed
			store.Update(ctx, inc)

			log.Println("Execution successful. Starting verification...")

			// ---------------- 3️⃣ VERIFY ----------------
			time.Sleep(30 * time.Second) // stabilization wait

			newCPU, err := aws.GetCPUUtilization(cfg, inc.InstanceID)
			if err != nil {
				log.Println("Verification CPU fetch failed:", err)
				continue
			}

			// Store verification details
			inc.CPUAfter = newCPU
			inc.ExecutionTimeSec = int(time.Since(startTime).Seconds())

			if newCPU < inc.CPU {
				now := time.Now()
				inc.State = incident.Verified
				inc.VerifiedAt = &now

				log.Println("Verification successful. CPU improved.")
			} else {
				log.Println("CPU did not improve. Rolling back...")

				err := scaler.Rollback(ctx, inc.InstanceID, inc.InstanceType)
				if err != nil {
					log.Println("Rollback failed:", err)
				}

				inc.State = incident.RolledBack
			}

			err = store.Update(ctx, inc)
			if err != nil {
				log.Println("Failed to update incident:", err)
			}

			// ---------------- DELETE MESSAGE ----------------
			_, err = sqsClient.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})

			if err != nil {
				log.Println("Delete message error:", err)
			}

			log.Println("✅ Incident lifecycle completed.")
		}
	}
}