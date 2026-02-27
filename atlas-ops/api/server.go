package api 

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"atlas-ops/aws"
	"github.com/aws/aws-sdk-go-v2/config"

)

func StartServer() {
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/monitor", monitorHandler)

	log.Println(" Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func monitorHandler(w http.ResponseWriter, r *http.Request) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		http.Error(w, "AWS config error", http.StatusInternalServerError)
		return
	}

instanceID, err := aws.GetFirstInstanceID(cfg)
	if err != nil {
		http.Error(w, "No EC2 instances found", http.StatusInternalServerError)
		return
	}

	instanceType, err := aws.GetInstanceType(cfg, instanceID)
	if err != nil {
		http.Error(w, "Failed to get instance type", http.StatusInternalServerError)
		return
	}

	cpu, err := aws.GetCPUUtilization(cfg, instanceID)
	if err != nil {
		http.Error(w, "Failed to fetch CPU metrics", http.StatusInternalServerError)
		return
	}

	result := policy.EvaluatePolicy(cpu, instanceType)

	response := map[string]interface{}{
		"instance_id":    instanceID,
		"instanceType":   instanceType,
		"cpu":            result.CPU,
		"status":         result.Status,
		"recommendation": result.Recommendation,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}