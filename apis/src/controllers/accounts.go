package controllers

import (
	dynamodbpkg "GoldChain/apis/src/apis/src/DynamoDB"
	constantpkg "GoldChain/apis/src/apis/src/constant"
	pinpointpkg "GoldChain/apis/src/apis/src/pinpoint"
	utilspkg "GoldChain/apis/src/apis/src/utils"
	errorservice "assistant/ErrorService"
	"fmt"
	"net/mail"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetAccountInfo(c *gin.Context) {

	userName := c.Query("userName")

	resp, err := dynamodbpkg.GetDataFromDynamoDb("username", userName, constantpkg.D_ACCOUNT_TABLE)
	fmt.Println("response from dynamodb is ", resp)
	if err != nil {
		errorservice.ErrorResponse(c, 500, "Internal server error")
	}

	flag := resp["keyexist"].(bool)
	if !flag {
		errorservice.ErrorResponse(c, 400, "User doesn't exist.Please enter valid username")
	}
	respData := resp["data"].(map[string]interface{})
	delete(respData, "password")
	c.JSON(200, respData)
}

func CreateUserAccount(c *gin.Context) {

	reqBody, _ := utilspkg.GetReqBodyMap(c)
	var requiredFileds = []string{"username", "password", "name"}
	flag := utilspkg.AttributesExistInMap(reqBody, requiredFileds)
	if !flag {
		errorservice.ErrorResponse(c, 500, "Invalid request body")
		return
	}
	userName := reqBody["username"].(string)
	passWord := reqBody["password"].(string)
	name := reqBody["name"].(string)
	_, err := mail.ParseAddress(userName)
	if err != nil {
		errorservice.ErrorResponse(c, 500, "Please enter valid email address")
		return
	}

	result, err := dynamodbpkg.GetDataFromDynamoDb("username", userName, constantpkg.D_ACCOUNT_TABLE)
	if err != nil {
		errorservice.ErrorResponse(c, 500, "Internal Server Error")
		return
	}
	if result["keyexist"].(bool) == true {
		// errorservice.ErrorResponse(c, 400, "user already exist")
		// return
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(passWord), bcrypt.DefaultCost)
	dynamoDbData := map[string]*dynamodb.AttributeValue{
		"username": {
			S: aws.String(userName),
		},
		"password": {
			S: aws.String(string(hashedPassword)),
		},
		"Name": {
			S: aws.String(name),
		},
	}
	_, err = dynamodbpkg.InsertDataInDynamoDBTable(constantpkg.D_ACCOUNT_TABLE, dynamoDbData)
	if err != nil {
		errorservice.ErrorResponse(c, 500, "Internal Server Error !!!. Please retry")
		return
	}
	response := map[string]interface{}{

		"message": "user is created",
		"status":  "SUCCESS",
	}

	toAddress := make([]*string, 0)
	toAddress = append(toAddress, aws.String(userName))
	flag, err = pinpointpkg.SendPinPointAccountCreationMail(userName, toAddress)
	if err != nil {
		errorservice.ErrorResponse(c, 500, err.Error())
		return
	}

	c.JSON(200, gin.H{"response ": response})
	return
}

func LoginAccount(c *gin.Context) {

	reqBody, err := utilspkg.GetReqBodyMap(c)
	if err != nil {
		errorservice.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	requiredField := []string{"username", "granttype", "password"}
	flag := utilspkg.AttributesExistInMap(reqBody, requiredField)
	if !flag {
		errorservice.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	grantType := reqBody["granttype"].(string)
	if strings.ToLower(grantType) != "password" {
		errorservice.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	passWord := reqBody["password"].(string)
	userName := reqBody["username"].(string)

	dbResponse, err := dynamodbpkg.GetDataFromDynamoDb("username", userName, constantpkg.D_ACCOUNT_TABLE)
	if err != nil {
		fmt.Println("Unable to get username data from dynamodb", err)
		errorservice.ErrorResponse(c, 500, "Internal server error")
		return
	}
	keyExist := dbResponse["keyexist"].(bool)
	if !keyExist {
		errorservice.ErrorResponse(c, 500, "user doesn't exist")
		return

	}
	dbPassWord := dbResponse["data"].(map[string]interface{})["password"].(string)
	err = bcrypt.CompareHashAndPassword([]byte(dbPassWord), []byte(passWord))
	if err != nil {
		errorservice.ErrorResponse(c, 500, "username or password is incorrect !!!! ")
		return
	}
	response := map[string]interface{}{
		"username": userName,
		"grant":    "allowed",
		"status":   "SUCCESS",
	}
	c.JSON(200, response)
	return

}

func UpdateAccount(c *gin.Context) {
	reqBody, err := utilspkg.GetReqBodyMap(c)
	if err != nil {
		errorservice.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	requiredField := []string{"userName", "address", "updateType"}
	flag := utilspkg.AttributesExistInMap(reqBody, requiredField)
	if !flag {
		errorservice.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	updateType := reqBody["updateType"].(string)
	if strings.ToLower(updateType) != "address" {
		errorservice.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	userName := reqBody["userName"].(string)
	addressMap := reqBody["address"].(map[string]interface{})
	requiredAddressFields := []string{"addressLine1", "addressLine2", "state", "city", "zipCode", "country"}
	flag = utilspkg.AttributesExistInMap(addressMap, requiredAddressFields)
	if !flag {
		errorservice.ErrorResponse(c, 400, "Please enter complete address")
	}
	exist, err := dynamodbpkg.GetDataFromDynamoDb("username", userName, constantpkg.D_ACCOUNT_TABLE)
	if err != nil {
		errorservice.ErrorResponse(c, 500, "Internal server error !!!! DB_Error")
		return
	}
	if !exist["keyexist"].(bool) {
		errorservice.ErrorResponse(c, 500, "user doesn't exist !!!! .Please enter valid username")
		return
	}
	updateExpression := aws.String("SET Address = :Address")
	dynamoDbData := map[string]*dynamodb.AttributeValue{
		":Address": {
			M: map[string]*dynamodb.AttributeValue{
				"AddressLine1": {
					S: aws.String(addressMap["addressLine1"].(string)),
				},
				"AddressLine2": {
					S: aws.String(addressMap["addressLine2"].(string)),
				},
				"State": {
					S: aws.String(addressMap["state"].(string)),
				},
				"City": {
					S: aws.String(addressMap["city"].(string)),
				},
				"ZipCode": {
					S: aws.String(addressMap["zipCode"].(string)),
				},
				"Country": {
					S: aws.String(addressMap["country"].(string)),
				},
			},
		},
	}

	key := map[string]*dynamodb.AttributeValue{
		"username": {
			S: aws.String(userName),
		},
	}
	_, err = dynamodbpkg.UpdateDyanamoDBRow(constantpkg.D_ACCOUNT_TABLE, updateExpression, key, dynamoDbData)
	if err != nil {
		fmt.Println("error from db in updating ", err)
		errorservice.ErrorResponse(c, 500, "Internal server error. Update DBError !!!")
		return
	}
	response := map[string]interface{}{
		"username": userName,
		"status":   "SUCCESS",
		"message":  "Address has been updated !!!!",
	}
	c.JSON(200, response)
	return
}

func DeleteAccount(c *gin.Context) {
	userName := c.Query("userName")
	resp, err := dynamodbpkg.GetDataFromDynamoDb("username", userName, constantpkg.D_ACCOUNT_TABLE)
	if err != nil {
		errorservice.ErrorResponse(c, 500, "Internal server error")
		return
	}
	exist := resp["keyexist"].(bool)
	if !exist {
		errorservice.ErrorResponse(c, 400, "User does't exist. Please enter valid user !!")
		return
	}
	input := dynamodb.DeleteItemInput{
		TableName: aws.String(constantpkg.D_ACCOUNT_TABLE),
		Key: map[string]*dynamodb.AttributeValue{
			"username": {
				S: aws.String(userName),
			},
		},
	}
	flag, err := dynamodbpkg.DeleteDynamoDBRow(input)
	if err != nil || flag == false {
		errorservice.ErrorResponse(c, 500, "Internal server error")
		return
	}
	response := map[string]interface{}{
		"status":   "SUCCESS",
		"message":  "User is deleted",
		"username": userName,
	}
	c.JSON(200, response)
	return
}
