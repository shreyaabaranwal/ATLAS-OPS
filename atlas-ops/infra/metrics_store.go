package infra

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Store struct {
	Client *dynamodb.Client
	Table  string
}

type CPUMetric struct {
	InstanceID string  `dynamodbav:"instance_id"`
	Timestamp  int64   `dynamodbav:"timestamp"`
	CPU        float64 `dynamodbav:"cpu"`
}

func NewStore(client *dynamodb.Client) *Store {
	return &Store{
		Client: client,
		Table:  "atlas-metrics",
	}
}

// Save CPU snapshot
func (s *Store) SaveMetric(ctx context.Context, instanceID string, cpu float64) error {

	m := CPUMetric{
		InstanceID: instanceID,
		Timestamp:  time.Now().Unix(),
		CPU:        cpu,
	}

	item, err := attributevalue.MarshalMap(m)
	if err != nil {
		return err
	}

	_, err = s.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(s.Table),
		Item:      item,
	})

	return err
}

// Fetch last 5 minutes of metrics
func (s *Store) GetLast5Minutes(ctx context.Context, instanceID string) ([]CPUMetric, error) {

	now := time.Now().Unix()
	start := now - 300

	out, err := s.Client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(s.Table),
		KeyConditionExpression: aws.String("instance_id = :id AND #ts > :t"),
		ExpressionAttributeNames: map[string]string{
			"#ts": "timestamp",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: instanceID},
			":t":  &types.AttributeValueMemberN{Value: strconv.FormatInt(start, 10)},
		},
	})

	if err != nil {
		return nil, err
	}

	var metrics []CPUMetric
	err = attributevalue.UnmarshalListOfMaps(out.Items, &metrics)

	return metrics, err
}