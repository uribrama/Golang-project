package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambdaProxy *ginadapter.GinLambda

func main() {
	router := initGinRouter()
	env := os.Getenv("GO_ENV")
	if env == "development" {
		router.Run(":8080")
	} else {
		ginLambdaProxy = ginadapter.New(router)
		lambda.Start(Handler)
	}
}

func initGinRouter() *gin.Engine {
	handler := PersonRequestHandlerInit()
	gin := gin.Default()

	gin.Use(handler.configRequest)

	gin.GET("/persons/healthcheck", handler.healthcheck)

	routerGroup := gin.Group("/persons")
	routerGroup.Use(handler.authHandler)
	routerGroup.GET("/", handler.getPersons)
	routerGroup.PATCH("/:id", handler.updatePerson)
	routerGroup.POST("/", handler.createPerson)
	routerGroup.POST("/batch", handler.batchCSV)
	return gin
}

func Handler(cx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambdaProxy.ProxyWithContext(cx, request)
}
