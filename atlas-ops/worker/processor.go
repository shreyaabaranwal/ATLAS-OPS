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

			// 🔥 1️⃣ Atomic Lock: APPROVED → EXECUTING
			err := store.TransitionState(ctx, incidentID, incident.Approved, incident.Executing)
			if err != nil {
				log.Println("Another worker already processing or invalid state. Skipping.")
				continue
			}

			// 🔥 Increment execution attempts
			_ = store.IncrementAttempts(ctx, incidentID)

			// Fetch updated record
			inc, err := store.Get(ctx, incidentID)
			if err != nil {
				log.Println("Incident fetch error:", err)
				continue
			}

			targetType := "t3.medium"

			// ---------------- DRY RUN ----------------
			err = scaler.DryRun(ctx, inc.InstanceID, targetType)
			if err != nil {

				store.TransitionState(ctx, inc.ID, incident.Executing, incident.RolledBack)

				_ = audit.Save(ctx, &incident.AuditLog{
					IncidentID: inc.ID,
					Action:     "DRY_RUN",
					Actor:      "SYSTEM",
					Result:     "FAILED",
					CPUBefore:  inc.CPU,
				})

				continue
			}

			_ = audit.Save(ctx, &incident.AuditLog{
				IncidentID: inc.ID,
				Action:     "DRY_RUN",
				Actor:      "SYSTEM",
				Result:     "SUCCESS",
				CPUBefore:  inc.CPU,
			})

			// ---------------- EXECUTE ----------------
			startTime := time.Now()

			err = scaler.Execute(ctx, inc.InstanceID, targetType)
			if err != nil {

				store.TransitionState(ctx, inc.ID, incident.Executing, incident.RolledBack)

				_ = audit.Save(ctx, &incident.AuditLog{
					IncidentID: inc.ID,
					Action:     "EXECUTION",
					Actor:      "SYSTEM",
					Result:     "FAILED",
					CPUBefore:  inc.CPU,
				})

				continue
			}

			_ = store.TransitionState(ctx, inc.ID, incident.Executing, incident.Executed)

			_ = audit.Save(ctx, &incident.AuditLog{
				IncidentID: inc.ID,
				Action:     "EXECUTION",
				Actor:      "SYSTEM",
				Result:     "EXECUTED",
				CPUBefore:  inc.CPU,
			})

			// ---------------- VERIFY ----------------
			time.Sleep(30 * time.Second)

			newCPU, err := aws.GetCPUUtilization(cfg, inc.InstanceID)
			if err != nil {
				continue
			}

			inc.CPUAfter = newCPU
			inc.ExecutionTimeSec = int(time.Since(startTime).Seconds())

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

				_ = audit.Save(ctx, &incident.AuditLog{
					IncidentID: inc.ID,
					Action:     "ROLLBACK",
					Actor:      "SYSTEM",
					Result:     "ROLLED_BACK",
					CPUBefore:  inc.CPU,
					CPUAfter:   newCPU,
				})
			}

			_, _ = sqsClient.Client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      &queueURL,
				ReceiptHandle: msg.ReceiptHandle,
			})

			log.Println("✅ Incident lifecycle completed safely.")
		}
	}
}