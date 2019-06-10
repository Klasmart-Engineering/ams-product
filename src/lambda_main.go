// +build lambda

package main

import (
	"bitbucket.org/calmisland/product-lambda-funcs/src/server"
	"bitbucket.org/calmisland/go-server-shared/services/aws/awslambda"
)

func main() {
	server.Setup()
	initLambdaFunctions()
	awslambda.StartRestHandler(rootRouter)
}
