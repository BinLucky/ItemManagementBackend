package main

import (
	"os"

	"github.com/BinLucky/ItemManagementBackend/pkg/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return
	}
	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)

}

const tableName = "LambdaAwsInGoItemManagement"

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {

	switch req.HTTPMethod {
	case "GET":
		return handlers.GetItem(req, tableName, dynaClient)
	case "POST":
		return handlers.CreateItem(req, tableName, dynaClient)
	case "PUT":
		return handlers.UpdateItem(req, tableName, dynaClient)
	case "DELETE":
		return handlers.DeleteItem(req, tableName, dynaClient)
	default:
		return handlers.UnHandledMethod()
	}

}
