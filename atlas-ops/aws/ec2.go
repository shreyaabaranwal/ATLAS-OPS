package aws
import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

)

func ListEC2Instances(cfg aws.Config) error {
	client := ec2.NewFromConfig(cfg)

	result, err := client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil{
		return fmt.Errorf("failed to describe instance: %w", err)

	}
	if len(result.Reservations) == 0 {
		fmt.Println("No EC2 Instances found.")
		return nil
	}

	fmt.Println("=========EC2 INSTANCES =========")

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			fmt.Println("---------------------------------")

			if instance.InstanceId != nil {
				fmt.Println("Instance ID :", *instance.InstanceId)
			}
			fmt.Println("Insatnce Type :", instance.InstanceType)

			if instance.State != nil {
				fmt.Println("State        :", instance.State.Name)

			}

			if instance.PublicIpAddress != nil {
				fmt.Println("Public IP  :", *instance.PublicIpAddress)

			} else{
				fmt.Print("Public IP     :, N/A")

			}

           if instance.PrivateIpAddress != nil {
			fmt.Println("Private IP    :", *instance.PrivateIpAddress)

		   }
		   fmt.Println("-------------------------")
		   fmt.Println()
		}
	}


	return nil 
}


// GetFirstInstanceID returns first EC2 instance ID
func GetFirstInstanceID(cfg aws.Config) (string, error) {

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

	return "", fmt.Errorf("no instances found")
}