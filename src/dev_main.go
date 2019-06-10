// +build !lambda

package main

import (
	"bitbucket.org/calmisland/product-lambda-funcs/src/server"
	"bitbucket.org/calmisland/go-server-shared/servers/restserver"
)

func main() {
	server.Setup()
	initLambdaFunctions()
	initLambdaDevFunctions()

	restServer := &restserver.Server{
		ListenAddress: ":8044",
		Handler:       rootRouter,
	}

	err := restServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
