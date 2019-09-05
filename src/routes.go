package main

import (
	"bitbucket.org/calmisland/go-server-auth/authmiddlewares"
	"bitbucket.org/calmisland/go-server-requests/apirouter"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/src/handlers"
)

var (
	rootRouter *apirouter.Router
)

func initLambdaFunctions() {
	rootRouter = &apirouter.Router{}
	routerV1 := createLambdaRouterV1()
	rootRouter.AddRouter("v1", routerV1)
}

func createLambdaRouterV1() *apirouter.Router {
	router := apirouter.NewRouter()

	router.AddMiddleware(authmiddlewares.ValidateSession(globals.AccessTokenValidator, true))

	router.AddMethodHandler("GET", "content", handlers.HandleContentInfoMultiple)
	router.AddMethodHandler("GET", "product", handlers.HandleProductInfoListByIds)
	router.AddMethodHandler("GET", "serverinfo", handlers.HandleServerInfo)

	contentRouter := apirouter.NewRouter()
	contentRouter.AddMethodHandlerWildcard("GET", "contentId", handlers.HandleContentInfo)
	router.AddRouter("content", contentRouter)

	specificContentRouter := apirouter.NewRouter()
	specificContentRouter.AddMethodHandler("GET", "icon", handlers.HandleContentIconDownload)
	contentRouter.AddRouterWildcard("productId", specificContentRouter)

	productRouter := apirouter.NewRouter()
	productRouter.AddMethodHandlerWildcard("GET", "productId", handlers.HandleProductInfo)
	router.AddRouter("product", productRouter)

	specificProductRouter := apirouter.NewRouter()
	specificProductRouter.AddMethodHandler("GET", "icon", handlers.HandleProductIconDownload)
	productRouter.AddRouterWildcard("productId", specificProductRouter)

	return router
}
