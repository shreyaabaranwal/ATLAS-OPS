package api

import (
	"encoding/json"
	"log"
	"net/http"

	"atlas-ops/aws"
	"atlas-ops/policy"

	sdkaws "github.com/aws/aws-sdk-go-v2/aws"
)

func StartServer(cfg sdkaws.Config) {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

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

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}