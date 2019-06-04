// +build lambda

package main

import (
	"context"

	"bitbucket.org/calmisland/product-lambda-funcs/src/server"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

func handleRequest(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	lambdaCtx, _ := lambdacontext.FromContext(ctx)
	resp, err := handleLambdaRequest(lambdaCtx, &req)
	return *resp, err
}

func main() {
	server.Setup()
	initLambdaFunctions()
	lambda.Start(handleRequest)
}
