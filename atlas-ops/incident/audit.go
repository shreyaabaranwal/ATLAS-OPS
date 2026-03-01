package incident

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type AuditLog struct {
	IncidentID string  `dynamodbav:"incident_id"`
	Timestamp  int64   `dynamodbav:"timestamp"`
	Action     string  `dynamodbav:"action"`
	Actor      string  `dynamodbav:"actor"`
	Result     string  `dynamodbav:"result"`
	CPUBefore  float64 `dynamodbav:"cpu_before,omitempty"`
	CPUAfter   float64 `dynamodbav:"cpu_after,omitempty"`
}

type AuditStore struct {
	Client *dynamodb.Client
	Table  string
}

func NewAuditStore(cfg aws.Config, table string) *AuditStore {
	return &AuditStore{
		Client: dynamodb.NewFromConfig(cfg),
		Table:  table,
	}
}

func (a *AuditStore) Save(ctx context.Context, logEntry *AuditLog) error {

	if logEntry.Timestamp == 0 {
		logEntry.Timestamp = time.Now().Unix()
	}

	item, err := attributevalue.MarshalMap(logEntry)
	if err != nil {
		return fmt.Errorf("marshal audit failed: %w", err)
	}

	_, err = a.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(a.Table),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("audit put failed: %w", err)
	}

	return nil
}

func (a *AuditStore) GetTimeline(ctx context.Context, incidentID string) ([]AuditLog, error) {

	out, err := a.Client.Query(ctx, &dynamodb.QueryInput{
		TableName: aws.String(a.Table),
		KeyConditionExpression: aws.String("incident_id = :id"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":id": &types.AttributeValueMemberS{Value: incidentID},
		},
		ScanIndexForward: aws.Bool(true),
	})

	if err != nil {
		return nil, err
	}

	var logs []AuditLog
	err = attributevalue.UnmarshalListOfMaps(out.Items, &logs)

	return logs, err
}