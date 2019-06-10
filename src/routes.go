package main

import (
	"bitbucket.org/calmisland/go-server-shared/requests/apirequests"
	"bitbucket.org/calmisland/product-lambda-funcs/src/handlers"
)

var (
	rootRouter *apirequests.Router
)

func initLambdaFunctions() {
	rootRouter = &apirequests.Router{}
	routerV1 := createLambdaRouterV1()
	rootRouter.AddRouter("v1", routerV1)
}

func createLambdaRouterV1() *apirequests.Router {
	router := &apirequests.Router{}
	router.AddMethodHandler("GET", "product", handlers.HandleProductInfoListByIds)
	router.AddMethodHandler("GET", "serverinfo", handlers.HandleServerInfo)

	productRouter := &apirequests.Router{}
	productRouter.AddMethodHandlerWildcard("GET", "productId", handlers.HandleProductInfo)
	router.AddRouter("product", productRouter)

	return router
}
