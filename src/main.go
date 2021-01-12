// +build !lambda

package main

import (
	"context"

	"bitbucket.org/calmisland/go-server-configs/configs"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-requests/apirouter"
	"bitbucket.org/calmisland/go-server-requests/apiservers/httpserver"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/src/handlers"
	"bitbucket.org/calmisland/product-lambda-funcs/src/setup/globalsetup"
)

func main() {
	err := configs.UpdateConfigDirectoryPath(configs.DefaultConfigFolderName)
	if err != nil {
		panic(err)
	}

	globalsetup.Setup()
	rootRouter := handlers.InitializeRoutes()

	initLambdaDevFunctions(rootRouter)

	restServer := &httpserver.Server{
		ListenAddress: ":8044",
		Handler:       rootRouter,
	}

	err = restServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func initLambdaDevFunctions(rootRouter *apirouter.Router) {
	devRouter := apirouter.NewRouter()
	devRouter.AddMethodHandler("GET", "createtables", createTablesRequest)
	rootRouter.AddRouter("dev", devRouter)
}

func createTablesRequest(_ context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	err := globals.ProductDatabase.CreateDatabaseTables()
	if err != nil {
		return resp.SetServerError(err)
	}
	return nil
}
