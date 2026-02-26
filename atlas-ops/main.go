package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
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

fmt.Println("EC2 instances:")

for _, reservation := range result.Reservations {
for +, instance  := range reservation.Instances {
	fmt.Println("Instance ID:", *instance.InstanceId)
	fmt.Println("Instance Type:", instance.InstanceType)
	fmt.Println("State:", instance.State.Name)
	fmt.Println("Public IP:", *instance.PublicIpAddress)
	fmt.Println("Private IP:", *instance.PrivateIpAddress)
	fmt.Println("Launch Time:", instance.LaunchTime)
	fmt.Println("Tags:")
	for _, tag := range instance.Tags {
		fmt.Printf("  %s: %s\n", *tag.Key, *tag.Value)
	}
	fmt.Println()
}
}


}