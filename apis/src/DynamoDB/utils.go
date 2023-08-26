package dynamodbpkg

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var sess = session.Must(session.NewSessionWithOptions(session.Options{
	SharedConfigState: session.SharedConfigEnable,
}))
var svc = dynamodb.New(sess)

func GetDataFromDynamoDb(colname string, keyname string, tablename string) (map[string]interface{}, error) {
	input := dynamodb.GetItemInput{
		TableName: &tablename,
		Key: map[string]*dynamodb.AttributeValue{
			colname: {
				S: aws.String(keyname),
			},
		},
	}
	result, err := svc.GetItem(&input)
	if err != nil {
		fmt.Println("error is ", err.Error())
		return nil, err
	}
	var resultMap map[string]interface{}
	response := make(map[string]interface{})
	if result.Item == nil {
		response["keyexist"] = false
		response["data"] = nil
		return response, nil
	}
	dynamodbattribute.UnmarshalMap(result.Item, &resultMap)
	response["keyexist"] = true
	response["data"] = resultMap
	return response, nil
}

func InsertDataInDynamoDBTable(tablename string, data map[string]*dynamodb.AttributeValue) (bool, error) {
	input := dynamodb.PutItemInput{
		TableName: aws.String(tablename),
		Item:      data,
	}
	result, err := svc.PutItem(&input)
	if err != nil {
		return false, err
	}
	fmt.Println("response for PUT item from dynamodb ", result.String())
	return true, nil
}

func UpdateDyanamoDBRow(tableName string, updateExpression *string, key map[string]*dynamodb.AttributeValue, data map[string]*dynamodb.AttributeValue) (bool, error) {
	input := dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       key,
		ExpressionAttributeValues: data,
		UpdateExpression:          updateExpression,
	}
	_, err := svc.UpdateItem(&input)
	if err != nil {
		return false, err
	}
	return true, nil
}

func DeleteDynamoDBRow(input dynamodb.DeleteItemInput) (bool, error) {
	_, err := svc.DeleteItem(&input)
	if err != nil {
		return false, err
	}
	return true, nil

}
