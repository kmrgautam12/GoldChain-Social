package main

import (
	"GoldChain/apis/src/apis/src/controllers"
	"GoldChain/apis/src/apis/src/middleware"
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func main() {
	lambda.Start(handler)
	fmt.Println("Inside goldchain")

}
func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Println("inside handler function")
	router := gin.Default()
	ginLambda = ginadapter.New(router)
	router.POST("/token", controllers.GetGoldToken)
	router.Use(middleware.AuthMiddleWare)
	router.GET("account/info", controllers.GetAccountInfo)
	router.POST("/account/create", controllers.CreateUserAccount)
	router.PUT("/account/update", controllers.UpdateAccount)
	router.POST("/account/login", controllers.LoginAccount)
	router.DELETE("/account/delete", controllers.DeleteAccount)
	// resturant on boarding
	router.POST("/partner/create", controllers.CreateResturant)

	router.POST("/cart/add", controllers.AddItemToCart)
	return ginLambda.ProxyWithContext(ctx, request)
}
