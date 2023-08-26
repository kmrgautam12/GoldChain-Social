package utilspkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

func GetReqBodyMap(c *gin.Context) (map[string]interface{}, error) {
	var reqBodyMap map[string]interface{}
	reqBody := c.Request.Body
	b, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return nil, errors.New("Invalid request body")
	}
	err = json.Unmarshal(b, &reqBodyMap)
	if err != nil {
		return nil, errors.New("Invalid request body")
	}
	return reqBodyMap, nil

}

func AttributesExistInMap(m map[string]interface{}, a []string) bool {
	flag := true
	fmt.Println("request body is ", m)
	for _, attribute := range a {
		check := false
		for key, _ := range m {

			if key == attribute {
				check = true
			}
		}
		if !check {
			flag = false
			break
		}
	}
	return flag
}
