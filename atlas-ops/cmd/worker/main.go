package main

import (
	"context"
	"log"

	"atlas-ops/incident"
	"atlas-ops/worker"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {

	log.Println(" Starting Worker...")

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal("AWS config load failed:", err)
	}
	queueURL := "https://sqs.ap-south-1.amazonaws.com/606639293354/atlas-incident-queue"
	incidentStore := incident.NewDynamoStore(cfg, "atlas-incidents")
	auditStore := incident.NewAuditStore(cfg, "atlas-audit")

	worker.StartWorker(cfg, queueURL, incidentStore, auditStore)
}
