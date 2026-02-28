package execution

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
)

type EC2Scaler struct {
	Client *ec2.Client
}

func NewEC2Scaler(cfg aws.Config) *EC2Scaler {
	return &EC2Scaler{
		Client: ec2.NewFromConfig(cfg),
	}
}

// ---------------- DRY RUN ----------------
func (e *EC2Scaler) DryRun(ctx context.Context, instanceID string, newType string) error {

	log.Println("🔍 Running DryRun simulation...")

	input := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceID),
		InstanceType: &ec2types.AttributeValue{
			Value: aws.String(newType),
		},
		DryRun: aws.Bool(true),
	}

	_, err := e.Client.ModifyInstanceAttribute(ctx, input)

	if err != nil {

		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {

			// This means DryRun is SUCCESS
			if apiErr.ErrorCode() == "DryRunOperation" {
				log.Println("✅ DryRun permission validated.")
				return nil
			}
		}

		return fmt.Errorf("dry run failed: %w", err)
	}

	return nil
}

// ---------------- EXECUTE ----------------
func (e *EC2Scaler) Execute(ctx context.Context, instanceID string, newType string) error {

	log.Println("🚀 Executing scaling operation...")

	input := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceID),
		InstanceType: &ec2types.AttributeValue{
			Value: aws.String(newType),
		},
	}

	_, err := e.Client.ModifyInstanceAttribute(ctx, input)
	if err != nil {
		return fmt.Errorf("scaling failed: %w", err)
	}

	log.Println("✅ Scaling successful.")
	return nil
}

// ---------------- ROLLBACK ----------------
func (e *EC2Scaler) Rollback(ctx context.Context, instanceID string, previousType string) error {

	log.Println("↩ Rolling back instance type...")

	return e.Execute(ctx, instanceID, previousType)
}