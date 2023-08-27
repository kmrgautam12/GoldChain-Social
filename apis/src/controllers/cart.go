package controllers

import (
	dynamodbpkg "GoldChain/apis/src/DynamoDB"
	errorservice "GoldChain/apis/src/ErrorService"
	constantpkg "GoldChain/apis/src/constant"
	"GoldChain/apis/src/schema"
	utilspkg "GoldChain/apis/src/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func AddItemToCart(c *gin.Context) {
	reqBody, err := utilspkg.GetReqBodyMap(c)
	if err != nil {
		errorservice.ErrorResponse(c, 400, "Invalid request body")
		return
	}
	b, _ := json.Marshal(reqBody)
	var response schema.AddItemToCart
	err = json.Unmarshal(b, &response)
	if err != nil {

		errorservice.ErrorResponse(c, 400, err.Error())
		return
	}
	abandonedId := fmt.Sprintf("%v", time.Now().Unix()%1000000)
	taxAmt := fmt.Sprintf("%v", 0.18*response.Payload.Price)
	shipping := fmt.Sprintf("%v", 14.00)

	input := map[string]*dynamodb.AttributeValue{
		"Tax": {
			S: aws.String(taxAmt),
		},
		"Shipping": {
			S: aws.String(shipping),
		},
		"customer_id": {
			S: aws.String(response.ID),
		},
		"order_id": {
			S: aws.String((abandonedId)),
		},
		"payload": {
			M: map[string]*dynamodb.AttributeValue{
				"ItemId": {
					S: aws.String(response.Payload.ItemId),
				},
				"ItemName": {
					S: aws.String(response.Payload.ItemName),
				},

				"ResturantName": {
					S: aws.String(response.Payload.ResturantName),
				},
				"ResturantId": {
					S: aws.String(response.Payload.ResturantId),
				},
			},
		},
	}

	_, err = dynamodbpkg.InsertDataInDynamoDBTable(constantpkg.ABANDONED_ORDER_TABLE, input)
	if err != nil {
		errorservice.ErrorResponse(c, 500, "Internal Server Error")
		return
	}
	resp := map[string]interface{}{
		"customer_id":    response.ID,
		"product_price":  response.Payload.Price,
		"item_name":      response.Payload.ItemName,
		"resturant_name": response.Payload.ResturantName,
		"tax":            taxAmt,
		"shipping":       shipping,
		"response_status": map[string]interface{}{
			"message": "item added to cart",
			"status":  "SUCCESS",
			"code":    "200",
		},
	}
	c.JSON(200, resp)
}
