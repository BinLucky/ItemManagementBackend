package handlers

import (
	"net/http"

	"github.com/BinLucky/ItemManagement/pkg/item"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type ErrorBody struct {
	ErrorMsg *string `json:"error,omitempty`
}

var ErrorMethodNotAllowed = "Method not allowed"

func GetItem(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
	barcode := req.QueryStringParameters["barcode"]

	if len(barcode) > 0 {
		result, err := item.FetchItem(barcode, tableName, dynaClient)
		if err != nil {
			return apiResponse(http.StatusBadRequest, ErrorBody{aws.String(err.Error())})
		}
		return apiResponse(http.StatusOK, result)
	}

}

func CreateItem(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
}

func UpdateItem(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
}

func DeleteItem(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*events.APIGatewayProxyResponse, error) {
}

func UnHandledMethod() (*events.APIGatewayProxyResponse, error) {

	return apiResponse(http.StatusMethodNotAllowed, ErrorMethodNotAllowed)
}
