package main

import (
	"context"
	"log"


	"atlas-ops/incident"
	"atlas-ops/worker"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	queueURL := "https://sqs.ap-south-1.amazonaws.com/458329143405/atlas-incident-queue"

	// 🔥 FIX — pass table name
	tableName := "atlas-incidents"

	store := incident.NewDynamoStore(cfg, tableName)

	log.Println("👷 Starting Worker service...")
	worker.StartWorker(cfg, queueURL, store)
}