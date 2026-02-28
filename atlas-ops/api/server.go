package api

import (
	
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"atlas-ops/aws"
	"atlas-ops/execution"
	"atlas-ops/incident"
	"atlas-ops/policy"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
)

type JSONResponse struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func StartServer(cfg sdkaws.Config) {

	store := incident.NewDynamoStore(cfg)
	mux := http.NewServeMux()

	// ---------------- HEALTH ----------------
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		respondJSON(w, http.StatusOK, JSONResponse{Message: "Server healthy"})
	})

	// ---------------- MONITOR ----------------
	mux.HandleFunc("/monitor", func(w http.ResponseWriter, r *http.Request) {

		_, instanceType, cpu, err := fetchInstanceData(cfg)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Error: err.Error()})
			return
		}

		result := policy.EvaluatePolicy(cpu, instanceType)
		respondJSON(w, http.StatusOK, JSONResponse{Data: result})
	})

	// ---------------- CREATE INCIDENT ----------------
	mux.HandleFunc("/incident", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, JSONResponse{Error: "POST required"})
			return
		}

		instanceID, instanceType, cpu, err := fetchInstanceData(cfg)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Error: err.Error()})
			return
		}

		result := policy.EvaluatePolicy(cpu, instanceType)

		if result.Status != "HIGH_CPU" {
			respondJSON(w, http.StatusOK, JSONResponse{
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
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Error: err.Error()})
			return
		}

		logStateChange(inc.ID, string(inc.State))

		respondJSON(w, http.StatusCreated, JSONResponse{
			Message: "Incident created",
			Data:    map[string]string{"incident_id": inc.ID},
		})
	})

	// ---------------- APPROVE INCIDENT ----------------
	mux.HandleFunc("/approve/", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			respondJSON(w, http.StatusMethodNotAllowed, JSONResponse{Error: "POST required"})
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/approve/")

		inc, err := store.Get(r.Context(), id)
		if err != nil {
			respondJSON(w, http.StatusNotFound, JSONResponse{Error: "Incident not found"})
			return
		}

		scaler := execution.NewEC2Scaler(cfg)

		// APPROVED
		inc.State = incident.Approved
		store.Create(r.Context(), inc)
	logStateChange(id, string(inc.State))

		targetType := "t3.medium"

		// Dry Run
		if err := scaler.DryRun(inc.InstanceID, targetType); err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Error: err.Error()})
			return
		}

		// Execute
		if err := scaler.Execute(inc.InstanceID, targetType); err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Error: err.Error()})
			return
		}

		inc.State = incident.Executed
		store.Create(r.Context(), inc)
logStateChange(id, string(inc.State))

		time.Sleep(20 * time.Second)

		newCPU, err := aws.GetCPUUtilization(cfg, inc.InstanceID)
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, JSONResponse{Error: err.Error()})
			return
		}

		if newCPU > 80 {
			scaler.Rollback(inc.InstanceID, inc.InstanceType)
			inc.State = incident.Detected
			store.Create(r.Context(), inc)
			logStateChange(id, string(inc.State))

			respondJSON(w, http.StatusOK, JSONResponse{
				Message: "Rollback executed. Issue persists.",
			})
			return
		}

		inc.State = incident.Verified
		store.Create(r.Context(), inc)
	logStateChange(id, string(inc.State))

		respondJSON(w, http.StatusOK, JSONResponse{
			Message: "Incident resolved successfully",
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
	json.NewEncoder(w).Encode(payload)
}

func logStateChange(id string, state string) {
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