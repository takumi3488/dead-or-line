package main

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/takumi3488/dead-or-line/cmd/client"
)

func main() {
	// Load environment variables
	tableName, ok := os.LookupEnv("DYNAMODB_TABLE_NAME")
	if !ok {
		panic("DYNAMODB_TABLE_NAME is not set")
	}
	key, ok := os.LookupEnv("DYNAMODB_KEY")
	if !ok {
		panic("DYNAMODB_KEY is not set")
	}

	// Initialize DynamoDB client
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultRegion(""))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	dynamoClient := client.NewDynamoDB(dynamodb.NewFromConfig(cfg), tableName, key)
	dynamoClient.GetNotifiedAt(ctx)
}
