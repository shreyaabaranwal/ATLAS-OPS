package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func GetFirstEC2InstanceID(cfg aws.Config) (string, error) {
	client := ec2.NewFromConfig(cfg)

	result, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return "", fmt.Errorf("failed to describe instances: %w", err)
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.InstanceId != nil {
				return *instance.InstanceId, nil
			}
		}
	}

	return "", fmt.Errorf("no EC2 instances found")
}

func GetInstanceType(cfg aws.Config, instanceID string) (string, error) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to describe instance: %w", err)
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			if instance.InstanceType != "" {
				return string(instance.InstanceType), nil
			}
		}
	}

	return "", fmt.Errorf("instance type not found")
}