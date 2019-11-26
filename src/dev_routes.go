// +build !lambda

package main

import (
	"context"

	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-requests/apirouter"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
)

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
