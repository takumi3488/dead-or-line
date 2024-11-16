package main

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

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
	noticeInterval, ok := os.LookupEnv("NOTICE_INTERVAL")
	if !ok {
		panic("NOTICE_INTERVAL is not set")
	}
	noticeIntervalInt, err := strconv.Atoi(noticeInterval)
	if err != nil {
		panic("NOTICE_INTERVAL is not integer")
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

	// 経過時間がNOTICE_INTERVALを超えているか確認
	dynamoClient := client.NewDynamoDB(dynamodb.NewFromConfig(cfg), tableName, key)
	notifiedAt := dynamoClient.GetNotifiedAt(ctx)
	now := time.Now().Unix()
	if !(notifiedAt == 0 || now-notifiedAt > int64(noticeIntervalInt)) {
		return
	}

	// URLのステータスを確認
	ok, err = isOK(targetUrl)
	if err != nil {
		panic("URL check error, " + err.Error())
	}
	if ok {
		return
	}

	// Lineに通知
	message := lineClient.CreateMessage(targetUrl)
	lineClient.Notify(message)
}

func isOK(url string) (bool, error) {
	c := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	resp, err := c.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Redirect
	if resp.StatusCode >= 300 && resp.StatusCode <= 307 {
		return isOK(resp.Header.Get("Location"))
	}

	// Success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return true, nil
	}

	// Error
	return false, nil
}
