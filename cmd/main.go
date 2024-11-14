package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithDefaultRegion(""))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := dynamodb.NewFromConfig(cfg)

	input := &dynamodb.ListTablesInput{
		Limit: aws.Int32(5),
	}

	result, err := client.ListTables(ctx, input)
	if err != nil {
		panic("failed to list tables, " + err.Error())
	}

	for _, table := range result.TableNames {
		println(table)
	}
}
