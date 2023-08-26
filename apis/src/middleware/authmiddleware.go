package middleware

import (
	dynamodbpkg "GoldChain/apis/src/DynamoDB"
	constantpkg "GoldChain/apis/src/constant"

	errorservice "assistant/ErrorService"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare(c *gin.Context) {
	headers := c.Request.Header
	authorizarion := headers.Get("Authorization")
	if len(authorizarion) == 0 {
		errorservice.ErrorResponse(c, 401, "Authorization header can't be empty")
		c.Abort()
		return
	}

	resultMap, err := dynamodbpkg.GetDataFromDynamoDb("secret_name", "d_goldchain_token", constantpkg.D_SECRET_TABLE)
	if err != nil || !resultMap["keyexist"].(bool) {
		errorservice.ErrorResponse(c, 500, "Internal Server Error")
		c.Abort()
		return
	}
	secretVal := resultMap["data"].(map[string]interface{})["secret_value"].(string)

	token, err := jwt.ParseWithClaims(authorizarion, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretVal), nil
	})
	if err != nil {
		errorservice.ErrorResponse(c, 401, "unauthorized")
		c.Abort()
		return
	}
	if token.Valid {
		c.Next()
		return
	}
	errorservice.ErrorResponse(c, 401, "unauthorized")
	c.Abort()
}
