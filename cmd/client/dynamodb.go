package client

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDB struct {
	Client    *dynamodb.Client
	TableName string
	Key       string
}

func NewDynamoDB(client *dynamodb.Client, tableName string, key string) *DynamoDB {
	return &DynamoDB{
		Client:    client,
		TableName: tableName,
		Key:       key,
	}
}

func (d *DynamoDB) GetNotifiedAt(ctx context.Context) int64 {
	resp, err := d.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(d.TableName),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: d.Key,
			},
		},
	})
	if err != nil {
		slog.Error("failed to get item", "error", err)
		return math.MinInt64
	}
	notifiedAtString := resp.Item["NotifiedAt"].(*types.AttributeValueMemberN).Value
	notifiedAt, err := strconv.ParseInt(notifiedAtString, 10, 64)
	if err != nil {
		slog.Error("failed to parse NotifiedAt", "error", err)
		return math.MinInt64
	}
	return notifiedAt
}

func (d *DynamoDB) UpdateNotifiedAt(ctx context.Context, notifiedAt int64) {
	_, err := d.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(d.TableName),
		Key: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: d.Key,
			},
		},
		UpdateExpression: aws.String("SET NotifiedAt = :notifiedAt"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":notifiedAt": &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%d", notifiedAt),
			},
		},
	})
	if err != nil {
		slog.Error("failed to update item", "error", err)
	}
}
