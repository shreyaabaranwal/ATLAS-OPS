package main

import (
	"context"
	"fmt"
	"log"

	"atlas-ops/api"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	fmt.Println("AWS config loaded successfully")
	api.StartServer(cfg)
}