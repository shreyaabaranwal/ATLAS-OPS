package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

)

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil{
		log.Fatalf("unable to load SDK config, %v", err)

	}


client := ec2.NewFromConfig(cfg)

result, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})

if err != nil {
	log.Fatalf("unable to describe instances, %v", err)
}

fmt.Println("========== EC2 INSTANCES ==========")

if len(result.Reservations) == 0 {
	fmt.Println("No EC2 instances found.")
	return
}

for _, reservation := range result.Reservations {
	for _, instance := range reservation.Instances {

		fmt.Println("-------------------------------------------------")

		// Instance ID
		if instance.InstanceId != nil {
			fmt.Println("Instance ID      :", *instance.InstanceId)
		}

		// Instance Type
		fmt.Println("Instance Type    :", instance.InstanceType)

		// State
		if instance.State != nil {
			fmt.Println("State            :", instance.State.Name)
		}

		// Public IP
		if instance.PublicIpAddress != nil {
			fmt.Println("Public IP        :", *instance.PublicIpAddress)
		} else {
			fmt.Println("Public IP        : N/A")
		}

		// Private IP
		if instance.PrivateIpAddress != nil {
			fmt.Println("Private IP       :", *instance.PrivateIpAddress)
		} else {
			fmt.Println("Private IP       : N/A")
		}

		// Launch Time
		if instance.LaunchTime != nil {
			fmt.Println("Launch Time      :", instance.LaunchTime)
		}

		// Tags
		fmt.Println("Tags             :")
		if len(instance.Tags) == 0 {
			fmt.Println("  No Tags")
		} else {
			for _, tag := range instance.Tags {
				if tag.Key != nil && tag.Value != nil {
					fmt.Printf("  %s : %s\n", *tag.Key, *tag.Value)
				}
			}
		}

		fmt.Println("-------------------------------------------------")
		fmt.Println()
	}
}
}