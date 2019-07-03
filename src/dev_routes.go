// +build !lambda

package main

import (
	"bitbucket.org/calmisland/go-server-standard/databases/productdatabase"

	"bitbucket.org/calmisland/go-server-shared/v2/apierrors"
	"bitbucket.org/calmisland/go-server-shared/v2/requests/apirequests"
)

func initLambdaDevFunctions() {
	devRouter := &apirequests.Router{}
	devRouter.AddMethodHandler("GET", "createtables", createTablesRequest)
	rootRouter.AddRouter("dev", devRouter)
}

func createTablesRequest(req *apirequests.Request) (*apirequests.Response, error) {
	if req.HTTPMethod != "GET" {
		return apirequests.ClientError(apierrors.ErrorBadRequestMethod)
	}

	// Get the database
	db, err := productdatabase.GetDatabase()
	if err != nil {
		return apirequests.ServerError(err)
	}

	err = db.CreateDatabaseTables()
	if err != nil {
		return apirequests.ServerError(err)
	}

	return &apirequests.Response{
		StatusCode: 200,
		Body:       []byte("OK"),
	}, nil
}
