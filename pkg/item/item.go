package item

import (
	"encoding/json"
	"errors"

	"github.com/BinLucky/ItemManagementBackend/pkg/validators"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorFailedtoFetchRecord     = "Failed to fetch error"
	ErrorFailedToUnMarshalRecord = "Failed to UnMarshall recived record"
	ErorInvalidItemData          = "Invalid Item Data"
	ErrorInvalidItemBarcode      = "Invalid Barcode number"
	ErrorItemIsExist             = "The Item you want to create is already exist."
	ErrorFailedtoMarshalRecord   = "Failed to Marshall Record"
	ErrorCouldNotDynamoPutItem   = "The Item couldnt put to db"
	ErorItemIsNotExist           = "The Item which is you want to update is not exist"
	ErrorCouldNotDeleteItem      = "The Item could not delete which is you wish to delete."
)

type Item struct {
	Barcode  string `json:"barcode"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Location string `json:"location"`

	CurrentOwner string   `json:"currentowner"`
	OwnerHistory []string `json:"ownerhistory"`
}

func FetchItem(barcode, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Item, error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"barcode": {
				S: aws.String(barcode),
			},
		},
		TableName: aws.String(tableName),
	}
	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedtoFetchRecord)
	}
	resultItem := new(Item)
	err = dynamodbattribute.UnmarshalMap(result.Item, resultItem)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnMarshalRecord)
	}
	return resultItem, nil
}

func FetchItems(tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*[]Item, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}
	resultItems, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedtoFetchRecord)
	}
	incominItems := new([]Item)
	err = dynamodbattribute.UnmarshalListOfMaps(resultItems.Items, incominItems)

	if err != nil {
		return nil, errors.New(ErrorFailedToUnMarshalRecord)
	}
	return incominItems, nil
}

func CreateItem(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Item, error) {
	var cItem Item

	if err := json.Unmarshal([]byte(req.Body), &cItem); err != nil {
		return nil, errors.New(ErorInvalidItemData)
	}
	/* Validation */
	if !validators.IsBarcodeValid(cItem.Barcode) {
		return nil, errors.New(ErrorInvalidItemBarcode)
	}

	isItExist, _ := FetchItem(cItem.Barcode, tableName, dynaClient)

	if isItExist != nil && len(isItExist.Barcode) != 0 {
		return nil, errors.New(ErrorItemIsExist)
	}
	/* End of Validation */

	marshaledItem, err := dynamodbattribute.MarshalMap(cItem)
	if err != nil {
		return nil, errors.New(ErrorFailedtoMarshalRecord)
	}

	input := &dynamodb.PutItemInput{
		Item:      marshaledItem,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &cItem, nil

}

func DeleteItem(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) error {
	barcode := req.QueryStringParameters["barcode"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"barcode": {
				S: aws.String(barcode),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil

}

func UpdateItem(req events.APIGatewayProxyRequest, tableName string, dynaClient dynamodbiface.DynamoDBAPI) (*Item, error) {
	var uItem Item

	if err := json.Unmarshal([]byte(req.Body), &uItem); err != nil {

		return nil, errors.New(ErrorFailedToUnMarshalRecord)
	}

	if !validators.IsBarcodeValid(uItem.Barcode) {
		return nil, errors.New(ErrorInvalidItemBarcode)
	}

	isItExist, _ := FetchItem(uItem.Barcode, tableName, dynaClient)
	if isItExist == nil && len(isItExist.Barcode) == 0 {

		return nil, errors.New(ErorItemIsNotExist)

	}
	marshaledItem, err := dynamodbattribute.MarshalMap(uItem)
	if err != nil {
		return nil, errors.New(ErrorFailedtoMarshalRecord)
	}
	input := &dynamodb.PutItemInput{
		Item:      marshaledItem,
		TableName: &tableName,
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}
	return &uItem, nil

}
