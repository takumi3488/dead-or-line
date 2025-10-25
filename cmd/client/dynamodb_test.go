package client

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/testcontainers/testcontainers-go"
	dynamodule "github.com/testcontainers/testcontainers-go/modules/dynamodb"
)

const tableName = "table-name"
const key = "key"

func initDynamoDB(host string, port string) (*dynamodb.Client, error) {
	ctx := context.Background()
	// localhostの場合は明示的に127.0.0.1を使用してIPv6の問題を回避
	if host == "localhost" {
		host = "127.0.0.1"
	}
	cfg, err := config.LoadDefaultConfig(ctx, config.WithBaseEndpoint("http://"+host+":"+port))
	if err != nil {
		return nil, err
	}
	client := dynamodb.NewFromConfig(cfg)

	// DynamoDB Localが完全に起動するまで待機
	time.Sleep(2 * time.Second)

	_, err = client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []types.AttributeDefinition{{
			AttributeName: aws.String("Key"),
			AttributeType: types.ScalarAttributeTypeS,
		}},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Key"),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		return nil, err
	}

	_, err = client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item: map[string]types.AttributeValue{
			"Key": &types.AttributeValueMemberS{
				Value: key,
			},
			"NotifiedAt": &types.AttributeValueMemberN{
				Value: "0",
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func TestGetNotifiedAt(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "dummy")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "dummy")
	ctx := context.Background()
	ctr, err := dynamodule.Run(ctx, "amazon/dynamodb-local:2.5.3")
	defer func() {
		if err := testcontainers.TerminateContainer(ctr); err != nil {
			t.Fatal(err)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}
	host, err := ctr.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := ctr.MappedPort(ctx, "8000")
	if err != nil {
		t.Fatal(err)
	}
	client, err := initDynamoDB(host, port.Port())
	if err != nil {
		t.Fatal(err)
	}
	dynamodbClient := NewDynamoDB(client, tableName, key)
	notifiedAt := dynamodbClient.GetNotifiedAt(ctx)
	if notifiedAt != 0 {
		t.Errorf("expected 0, but got %d", notifiedAt)
	}
}

func TestUpdateNotifiedAt(t *testing.T) {
	t.Setenv("AWS_ACCESS_KEY_ID", "dummy")
	t.Setenv("AWS_SECRET_ACCESS_KEY", "dummy")
	ctx := context.Background()
	ctr, err := dynamodule.Run(ctx, "amazon/dynamodb-local:2.5.3")
	defer func() {
		if err := testcontainers.TerminateContainer(ctr); err != nil {
			t.Fatal(err)
		}
	}()
	if err != nil {
		t.Fatal(err)
	}
	host, err := ctr.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := ctr.MappedPort(ctx, "8000")
	if err != nil {
		t.Fatal(err)
	}
	client, err := initDynamoDB(host, port.Port())
	if err != nil {
		t.Fatal(err)
	}
	dynamodbClient := NewDynamoDB(client, tableName, key)
	dynamodbClient.UpdateNotifiedAt(ctx, 1)
	notifiedAt := dynamodbClient.GetNotifiedAt(ctx)
	if notifiedAt != 1 {
		t.Errorf("expected 1, but got %d", notifiedAt)
	}
}
