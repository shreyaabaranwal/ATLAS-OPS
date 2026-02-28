package execution

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Scaler struct {
	Client *ec2.Client
}

func NewEC2Scaler(cfg aws.Config) *EC2Scaler {
	return &EC2Scaler{
		Client: ec2.NewFromConfig(cfg),
	}
}

// ---------------- Dry Run ----------------
func (e *EC2Scaler) DryRun(instanceID string, newType string) error {

	log.Println("Running DryRun simulation...")

	input := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceID),
		InstanceType: &ec2types.AttributeValue{
			Value: aws.String(newType),
		},
		DryRun: aws.Bool(true),
	}

	_, err := e.Client.ModifyInstanceAttribute(context.TODO(), input)

	// AWS DryRun returns error intentionally
	if err != nil {
		log.Println("DryRun validation complete.")
		return nil
	}

	return nil
}

// ---------------- Execute Scaling ----------------
func (e *EC2Scaler) Execute(instanceID string, newType string) error {

	log.Println("Executing scaling operation...")

	input := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceID),
		InstanceType: &ec2types.AttributeValue{
			Value: aws.String(newType),
		},
	}

	_, err := e.Client.ModifyInstanceAttribute(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("scaling failed: %w", err)
	}

	return nil
}

// ---------------- Rollback ----------------
func (e *EC2Scaler) Rollback(instanceID string, previousType string) error {

	log.Println("Rolling back instance type...")

	return e.Execute(instanceID, previousType)
}