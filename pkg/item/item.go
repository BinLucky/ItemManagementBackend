package item

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var(
	ErrorFailedtoFetchRecord = "Failed to fetch error"
	ErrorFailedToUnMarshalRecord = "Failed to UnMarshall recived record"

)

type Item struct{
	Barcode string `json:"barcode"`
	Brand string   `json:"brand"`
	Model string	`json:"model"`
	Location string	`json:"location"`

	CurrentOwner string `json:"currentowner"`
	OwnerHistory []string	`json:"ownerhistory"`
}

func FetchItem(barcode, tableName string,dynaClient dynamodbiface.DynamoDBAPI )(*Item,error){

	input := &dynamodb.GetItemInput{
		Key : map[string]*dynamodb.AttributeValue{
			"barcode":{
				S: aws.String(barcode)
			}
		},
		TableName: aws.String(tableName),
	}
	result , err := dynaClient.GetItem(input)
	if err!=nil{
		return nil, errors.New(ErrorFailedtoFetchRecord)
	}
	resultItem := new(Item)
	err = dynamodbattribute.UnmarshalMap(result.Item,resultItem)
	if err!= nil{
		return nil, errors.New(ErrorFailedToUnMarshalRecord)
	}
	return resultItem , nil
}

func FetchItems()(){}

func CreateItem()(){}

func DeleteItem()(){}

func UpdateItem()(){}

