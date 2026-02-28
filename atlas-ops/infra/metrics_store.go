package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Metric struct {
	InstanceID string  `dynamodbav:"instance_id"`
	Timestamp  int64   `dynamodbav:"timestamp"`
	CPU        float64 `dynamodbav:"cpu"`
}

type MetricsStore struct {
	Client *dynamodb.Client
	Table  string
}

func NewMetricsStore(cfg aws.Config, table string) *MetricsStore {
	return &MetricsStore{
		Client: dynamodb.NewFromConfig(cfg),
		Table:  table,
	}
}

func (m *MetricsStore) SaveMetric(ctx context.Context, instanceID string, cpu float64) error {

	metric := Metric{
		InstanceID: instanceID,
		Timestamp:  time.Now().Unix(),
		CPU:        cpu,
	}

	item, err := attributevalue.MarshalMap(metric)
	if err != nil {
		return err
	}

	_, err = m.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(m.Table),
		Item:      item,
	})

	return err
}

func (m *MetricsStore) GetLast5Minutes(ctx context.Context, instanceID string) ([]Metric, error) {

	fiveMinAgo := time.Now().Add(-5 * time.Minute).Unix()

	out, err := m.Client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(m.Table),
		KeyConditionExpression: aws.String("instance_id = :id AND #ts > :time"),
		ExpressionAttributeNames: map[string]string{
			"#ts": "timestamp",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id":   &types.AttributeValueMemberS{Value: instanceID},
			":time": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", fiveMinAgo)},
		},
	})

	if err != nil {
		return nil, err
	}

	var metrics []Metric
	err = attributevalue.UnmarshalListOfMaps(out.Items, &metrics)

	return metrics, err
}