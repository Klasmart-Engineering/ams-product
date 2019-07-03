// +build lambda

package main

import (
	"bitbucket.org/calmisland/product-lambda-funcs/src/server"
	"bitbucket.org/calmisland/go-server-shared/v2/services/aws/awslambda"
	"bitbucket.org/calmisland/go-server-shared/v2/configs"
)

func main() {
	err := configs.UpdateConfigDirectoryPath(configs.DefaultConfigFolderName)
	if err != nil {
		panic(err)
	}
	server.Setup()
	initLambdaFunctions()
	awslambda.StartRestHandler(rootRouter)
}
