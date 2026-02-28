package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"atlas-ops/aws"
	"atlas-ops/incident"
	"atlas-ops/policy"
	"atlas-ops/queue"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
)

type JSONResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func StartServer(cfg sdkaws.Config) {

	store := incident.NewDynamoStore(cfg)

	// 🔥 Replace with your real Queue URL
	queueURL := "https://sqs.ap-south-1.amazonaws.com/458329143405/atlas-incident-queue"

	sqsClient := queue.NewSQSClient(cfg, queueURL)

	mux := http.NewServeMux()

	// ---------------- HEALTH ----------------
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, JSONResponse{Message: "Server healthy"})
	})

	// ---------------- MONITOR ----------------
	mux.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {

		_, instanceType, cpu, err := fetchInstanceData(cfg)
		if err != nil {
			respondJSON(w, 500, JSONResponse{Error: err.Error()})
			return
		}

		result := policy.EvaluatePolicy(cpu, instanceType)
		respondJSON(w, 200, JSONResponse{Data: result})
	})

	// ---------------- CREATE INCIDENT ----------------
	mux.HandleFunc("/incident", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			respondJSON(w, 405, JSONResponse{Error: "POST required"})
			return
		}

		instanceID, instanceType, cpu, err := fetchInstanceData(cfg)
		if err != nil {
			respondJSON(w, 500, JSONResponse{Error: err.Error()})
			return
		}

		result := policy.EvaluatePolicy(cpu, instanceType)

		if result.Status != "HIGH_CPU" {
			respondJSON(w, 200, JSONResponse{
				Message: "System healthy. No incident created.",
			})
			return
		}

		inc := &incident.Incident{
			ID:             generateID(),
			InstanceID:     instanceID,
			InstanceType:   instanceType,
			CPU:            cpu,
			State:          incident.Detected,
			Recommendation: result.Recommendation,
			CreatedAt:      time.Now(),
		}

		if err := store.Create(r.Context(), inc); err != nil {
			respondJSON(w, 500, JSONResponse{Error: err.Error()})
			return
		}

		logStateChange(inc.ID, inc.State)

		respondJSON(w, 201, JSONResponse{
			Message: "Incident created",
			Data:    map[string]string{"incident_id": inc.ID},
		})
	})

	// ---------------- APPROVE INCIDENT ----------------
	mux.HandleFunc("/approve/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			respondJSON(w, 405, JSONResponse{Error: "POST required"})
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/approve/")

		inc, err := store.Get(r.Context(), id)
		if err != nil {
			respondJSON(w, 404, JSONResponse{Error: "Incident not found"})
			return
		}

		// Mark Approved
		inc.State = incident.Approved
		if err := store.Create(r.Context(), inc); err != nil {
			respondJSON(w, 500, JSONResponse{Error: err.Error()})
			return
		}
		logStateChange(id, inc.State)

		// 🔥 Send to SQS instead of executing directly
		if err := sqsClient.SendMessage(r.Context(), inc.ID); err != nil {
			respondJSON(w, 500, JSONResponse{Error: err.Error()})
			return
		}

		respondJSON(w, 200, JSONResponse{
			Message: "Incident sent to processing queue",
		})
	})

	// ---------------- METRICS ----------------
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {

		metrics, err := store.GetMetrics(r.Context())
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{
				Error: err.Error(),
			})
			return
		}

		respondJSON(w, http.StatusOK, JSONResponse{
			Data: metrics,
		})
	})

	log.Println("🚀 Server running on :8081")
	log.Fatal(http.ListenAndServe(":8081", mux))
}

// ---------------- HELPERS ----------------

func generateID() string {
	return fmt.Sprintf("INC-%d", time.Now().UnixNano())
}

func respondJSON(w http.ResponseWriter, status int, payload JSONResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func logStateChange(id string, state incident.State) {
	log.Printf("[INCIDENT %s] State changed to %s", id, state)
}

func fetchInstanceData(cfg sdkaws.Config) (string, string, float64, error) {

	instanceID, err := aws.GetFirstEC2InstanceID(cfg)
	if err != nil {
		return "", "", 0, err
	}

	instanceType, err := aws.GetInstanceType(cfg, instanceID)
	if err != nil {
		return "", "", 0, err
	}

	cpu, err := aws.GetCPUUtilization(cfg, instanceID)
	if err != nil {
		return "", "", 0, err
	}

	return instanceID, instanceType, cpu, nil
}

// metrics handler moved inside StartServer; no standalone mux usage here