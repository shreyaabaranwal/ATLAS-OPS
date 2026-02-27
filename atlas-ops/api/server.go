package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"atlas-ops/aws"
	"atlas-ops/execution"
	"atlas-ops/incident"
	"atlas-ops/policy"
	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
)

func StartServer(cfg sdkaws.Config) {

	// Health Endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Monitoring Endpoint (optional)
	http.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {

		instanceID, err := aws.GetFirstEC2InstanceID(cfg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		instanceType, err := aws.GetInstanceType(cfg, instanceID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		cpu, err := aws.GetCPUUtilization(cfg, instanceID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		result := policy.EvaluatePolicy(cpu, instanceType)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	})

	// Create Incident Endpoint
	http.HandleFunc("/incident", func(w http.ResponseWriter, r *http.Request) {

		instanceID, err := aws.GetFirstEC2InstanceID(cfg)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		instanceType, err := aws.GetInstanceType(cfg, instanceID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		cpu, err := aws.GetCPUUtilization(cfg, instanceID)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		result := policy.EvaluatePolicy(cpu, instanceType)

		if result.Status != "HIGH_CPU" {
			w.Write([]byte("System healthy. No incident created."))
			return
		}

		inc := &incident.Incident{
			InstanceID:     instanceID,
			InstanceType:   instanceType,
			CPU:            cpu,
			State:          incident.Detected,
			Recommendation: result.Recommendation,
			CreatedAt:      time.Now(),
		}

		id := incident.Create(inc)

		log.Printf("[INCIDENT %s] State changed to %s", id, inc.State)

		inc.State = incident.Proposed
		log.Printf("[INCIDENT %s] State changed to %s", id, inc.State)

		w.Write([]byte("Incident created with ID: " + id))
	})

	// Approve Incident Endpoint
	http.HandleFunc("/approve/", func(w http.ResponseWriter, r *http.Request) {

		id := r.URL.Path[len("/approve/"):]

		inc, ok := incident.Get(id)
		if !ok {
			http.Error(w, "Incident not found", 404)
			return
		}

		scaler := execution.NewEC2Scaler(cfg)

		inc.State = incident.Approved
		log.Printf("[INCIDENT %s] State changed to %s", id, inc.State)

		// Dry Run
		err := scaler.DryRun(inc.InstanceID, "t3.medium")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		// Execute
		err = scaler.Execute(inc.InstanceID, "t3.medium")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		inc.State = incident.Executed
		log.Printf("[INCIDENT %s] State changed to %s", id, inc.State)

		inc.State = incident.Verified
		log.Printf("[INCIDENT %s] State changed to %s", id, inc.State)

		w.Write([]byte("Incident executed successfully"))
	})

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8081", nil))
}