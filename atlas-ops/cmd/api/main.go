package main

import (
	"context"
	"log"

	"atlas-ops/api"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	log.Println("🚀 Starting API service...")
	api.StartServer(cfg)
}