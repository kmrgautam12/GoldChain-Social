package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func init() {
	svc := dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{},
	}
}

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context) {

	// writing daily cron job of the create new account
}
func GetAllDynamoDbRecords() {

}
