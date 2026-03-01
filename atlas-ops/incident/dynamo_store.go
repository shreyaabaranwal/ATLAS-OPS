package incident

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoStore struct {
	Client *dynamodb.Client
	Table  string
}

// Table name should not be hardcoded in production
func NewDynamoStore(cfg aws.Config, tableName string) *DynamoStore {
	return &DynamoStore{
		Client: dynamodb.NewFromConfig(cfg),
		Table:  tableName,
	}
}

func (d *DynamoStore) Create(ctx context.Context, inc *Incident) error {

	item, err := attributevalue.MarshalMap(inc)
	if err != nil {
		return fmt.Errorf("marshal incident failed: %w", err)
	}

	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.Table),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("put item failed: %w", err)
	}

	return nil
}

func (d *DynamoStore) Get(ctx context.Context, id string) (*Incident, error) {

	out, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.Table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil, fmt.Errorf("get item failed: %w", err)
	}

	if out.Item == nil {
		return nil, errors.New("incident not found")
	}

	var inc Incident
	err = attributevalue.UnmarshalMap(out.Item, &inc)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %w", err)
	}

	return &inc, nil
}

func (d *DynamoStore) ScanAll(ctx context.Context) ([]*Incident, error) {

	var incidents []*Incident
	var lastEvaluatedKey map[string]types.AttributeValue

	for {
		out, err := d.Client.Scan(ctx, &dynamodb.ScanInput{
			TableName:         aws.String(d.Table),
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		for _, item := range out.Items {
			var inc Incident
			if err := attributevalue.UnmarshalMap(item, &inc); err != nil {
				continue
			}
			incidents = append(incidents, &inc)
		}

		if out.LastEvaluatedKey == nil {
			break
		}

		lastEvaluatedKey = out.LastEvaluatedKey
	}

	return incidents, nil
}
func (d *DynamoStore) Update(ctx context.Context, inc *Incident) error {

	item, err := attributevalue.MarshalMap(inc)
	if err != nil {
		return fmt.Errorf("marshal incident failed: %w", err)
	}

	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.Table),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("put item (update) failed: %w", err)
	}

	return nil
}



func (d *DynamoStore) TransitionState(
	ctx context.Context,
	id string,
	from State,
	to State,
) error {

	_, err := d.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(d.Table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("SET #s = :newState"),
		ConditionExpression: aws.String("#s = :expectedState"),
		ExpressionAttributeNames: map[string]string{
			"#s": "state",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":newState": &types.AttributeValueMemberS{
				Value: string(to),
			},
			":expectedState": &types.AttributeValueMemberS{
				Value: string(from),
			},
		},
	})

	return err
}

func (d *DynamoStore) IncrementAttempts(ctx context.Context, id string) error {

	_, err := d.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(d.Table),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression: aws.String("ADD execution_attempts :inc"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":inc": &types.AttributeValueMemberN{Value: "1"},
		},
	})

	return err
}