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
	targetUrl, ok := os.LookupEnv("TARGET_URL")
	if !ok {
		panic("TARGET_URL is not set")
	}
	tableName, ok := os.LookupEnv("DYNAMODB_TABLE_NAME")
	if !ok {
		panic("DYNAMODB_TABLE_NAME is not set")
	}
	key, ok := os.LookupEnv("DYNAMODB_KEY")
	if !ok {
		panic("DYNAMODB_KEY is not set")
	}
	lineToken, ok := os.LookupEnv("LINE_TOKEN")
	if !ok {
		panic("LINE_TOKEN is not set")
	}
	lineBaseMessage, ok := os.LookupEnv("LINE_BASE_MESSAGE")
	if !ok {
		panic("LINE_BASE_MESSAGE is not set")
	}
	lineTo, ok := os.LookupEnv("LINE_TO")
	if !ok {
		panic("LINE_TO is not set")
	}

	// Initialize DynamoDB client
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultRegion(""))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Initialize Line client
	lineClient := client.NewLine(lineToken, lineTo, lineBaseMessage)

	dynamoClient := client.NewDynamoDB(dynamodb.NewFromConfig(cfg), tableName, key)
	dynamoClient.GetNotifiedAt(ctx)

	message := lineClient.CreateMessage(targetUrl)
	lineClient.Notify(message)
}
