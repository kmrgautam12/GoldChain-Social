package controllers

import (
	dynamodbpkg "GoldChain/apis/src/apis/src/DynamoDB"
	constantpkg "GoldChain/apis/src/apis/src/constant"
	"GoldChain/apis/src/apis/src/schema"
	errorservice "assistant/ErrorService"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetGoldToken(c *gin.Context) {
	fmt.Println("remote address is ", c.Request.RemoteAddr)

	claimsStruct := &schema.Jwtclaims{

		Claims: jwt.StandardClaims{
			Id:        uuid.NewString(),
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
			Issuer:    "GoldChain-Team",
			IssuedAt:  time.Now().Unix(),
		},
	}

	secretsMap, err := dynamodbpkg.GetDataFromDynamoDb("secret_name", "d_goldchain_token", constantpkg.D_SECRET_TABLE)
	if err != nil {
		errorservice.ErrorResponse(c, 500, "Internal Server Error")
		return
	}
	fmt.Println("response from db is ", secretsMap)
	if secretsMap["keyexist"].(bool) == false {
		errorservice.ErrorResponse(c, 500, "Internal Server Error")
		return
	}

	secretValue := secretsMap["data"].(map[string]interface{})["secret_value"].(string)
	jwtStr := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsStruct.Claims)
	jwtToken, err := jwtStr.SignedString([]byte(secretValue))
	if err != nil {
		errorservice.ErrorResponse(c, 500, "error in signing the string")
		return
	}
	c.JSON(200, gin.H{"tokens": jwtToken})
}
