package api

import (
	"context"
	"log"
	"time"

	"atlas-ops/aws"
	"atlas-ops/execution"
	"atlas-ops/incident"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
)

var incidentQueue = make(chan string, 100)

// StartWorker starts background worker loop
func StartWorker(cfg sdkaws.Config, store *incident.DynamoStore) {

	go func() {
		log.Println("🚀 Background worker started")

		for id := range incidentQueue {
			processIncident(cfg, store, id)
		}
	}()
}

func processIncident(cfg sdkaws.Config, store *incident.DynamoStore, id string) {
	ctx := context.Background()

	log.Printf("[WORKER] Processing incident %s", id)

	inc, err := store.Get(ctx, id)
	if err != nil {
		log.Println("Worker error:", err)
		return
	}

	time.Sleep(20 * time.Second) // stabilization wait

	newCPU, err := aws.GetCPUUtilization(cfg, inc.InstanceID)
	if err != nil {
		log.Println("CPU fetch error:", err)
		return
	}

	scaler := execution.NewEC2Scaler(cfg)

	if newCPU > 80 {

	err := scaler.Rollback(ctx, inc.InstanceID, inc.InstanceType)
	if err != nil {
		log.Println("rollback failed:", err)
		return
	}

	inc.State = incident.RolledBack
	store.Update(ctx, inc)

	log.Println("Incident rolled back successfully")
	}

}