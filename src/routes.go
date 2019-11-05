package main

import (
	"bitbucket.org/calmisland/go-server-auth/authmiddlewares"
	"bitbucket.org/calmisland/go-server-requests/apirouter"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/src/handlers"
	"bitbucket.org/calmisland/product-lambda-funcs/src/middlewares"
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

	router.AddMethodHandler("GET", "serverinfo", handlers.HandleServerInfo)

	kidsLoopMiddlewareAuth := authmiddlewares.ValidateSession(globals.AccessTokenValidator, true)
	lpMiddlewareAuth := middlewares.ValidateSession(globals.AccessTokenValidator, globals.AccessTokenLPValidator, true)

	router.AddMethodHandler("GET", "content", handlers.HandleContentInfoMultiple, kidsLoopMiddlewareAuth)
	router.AddMethodHandler("GET", "product", handlers.HandleProductInfoListByIds, kidsLoopMiddlewareAuth)

	contentRouter := apirouter.NewRouter()
	contentRouter.AddMethodHandlerWildcard("GET", "contentId", handlers.HandleContentInfo, kidsLoopMiddlewareAuth)
	router.AddRouter("content", contentRouter)

	specificContentRouter := apirouter.NewRouter()
	specificContentRouter.AddMethodHandler("GET", "icon", handlers.HandleContentIconDownload, kidsLoopMiddlewareAuth)
	contentRouter.AddRouterWildcard("contentId", specificContentRouter)

	productRouter := apirouter.NewRouter()
	productRouter.AddMethodHandler("GET", "accesses", handlers.HandleAccessProductIdList, lpMiddlewareAuth)
	productRouter.AddMethodHandlerWildcard("GET", "productId", handlers.HandleProductInfo, kidsLoopMiddlewareAuth)
	productRouter.AddMethodHandlerWildcard("POST", "tickets", handlers.HandleTicketActivate, kidsLoopMiddlewareAuth)
	router.AddRouter("product", productRouter)

	specificProductRouter := apirouter.NewRouter()
	specificProductRouter.AddMethodHandler("GET", "icon", handlers.HandleProductIconDownload, kidsLoopMiddlewareAuth)
	productRouter.AddRouterWildcard("productId", specificProductRouter)

	return router
}
