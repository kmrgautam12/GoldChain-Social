package controllers

import (
	utilspkg "GoldChain/apis/src/apis/src/utils"
	errorservice "assistant/ErrorService"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func CreateResturant(c *gin.Context) {
	reqBody, err := utilspkg.GetReqBodyMap(c)
	if err != nil {
		errorservice.ErrorResponse(c, 400, "Invalid request body")
	}
	b, _ := json.Marshal(reqBody)
	var resturantMap map[string]interface{}
	json.Unmarshal(b, &resturantMap)

}
