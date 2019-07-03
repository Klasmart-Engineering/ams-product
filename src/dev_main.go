// +build !lambda

package main

import (
	"bitbucket.org/calmisland/go-server-shared/v2/configs"
	"bitbucket.org/calmisland/go-server-shared/v2/servers/restserver"
	"bitbucket.org/calmisland/product-lambda-funcs/src/server"
)

func main() {
	err := configs.UpdateConfigDirectoryPath(configs.DefaultConfigFolderName)
	if err != nil {
		panic(err)
	}
	server.Setup()
	initLambdaFunctions()
	initLambdaDevFunctions()

	restServer := &restserver.Server{
		ListenAddress: ":8044",
		Handler:       rootRouter,
	}

	err = restServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
