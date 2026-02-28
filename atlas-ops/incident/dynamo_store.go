package incident

import (
	"context"
	"fmt"
"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoStore struct {
	Client *dynamodb.Client
	Table  string
}

func NewDynamoStore(cfg aws.Config) *DynamoStore {
	return &DynamoStore{
		Client: dynamodb.NewFromConfig(cfg),
		Table:  "atlas-incidents",
	}
}

func (d *DynamoStore) Create(ctx context.Context, inc *Incident) error {

	item, err := attributevalue.MarshalMap(inc)
	if err != nil {
		return fmt.Errorf("failed to marshal incident: %w", err)
	}

	_, err = d.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(d.Table),
		Item:      item,
	})

	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
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
        return nil, fmt.Errorf("failed to get item: %w", err)
    }

    if out.Item == nil {
        return nil, fmt.Errorf("incident not found")
    }

    var inc Incident
    err = attributevalue.UnmarshalMap(out.Item, &inc)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal item: %w", err)
    }

    return &inc, nil
}

func (s *DynamoStore) ScanAll(ctx context.Context) ([]*Incident, error) {

	out, err := s.Client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(s.Table),
	})
	if err != nil {
		return nil, err
	}

	var incidents []*Incident

	for _, item := range out.Items {
		var inc Incident
		err := attributevalue.UnmarshalMap(item, &inc)
		if err != nil {
			continue
		}
		incidents = append(incidents, &inc)
	}

	return incidents, nil
}