package handlers

import (
	"bitbucket.org/calmisland/go-server-auth/authmiddlewares"
	"bitbucket.org/calmisland/go-server-requests/apirouter"
	"bitbucket.org/calmisland/go-server-requests/standardhandlers"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
)

var (
	rootRouter *apirouter.Router
)

// InitializeRoutes initializes the routes.
func InitializeRoutes() *apirouter.Router {
	if rootRouter != nil {
		return rootRouter
	}

	rootRouter = apirouter.NewRouter()
	if globals.CORSOptions != nil {
		rootRouter.AddCORSMiddleware(globals.CORSOptions)
	}

	routerV1 := createLambdaRouterV1()
	rootRouter.AddRouter("v1", routerV1)
	return rootRouter
}

func createLambdaRouterV1() *apirouter.Router {
	router := apirouter.NewRouter()
	router.AddMethodHandler("GET", "serverinfo", standardhandlers.HandleServerInfo)

	requireAuthMiddleware := authmiddlewares.ValidateSession(globals.AccessTokenValidator, true)

	router.AddMethodHandler("GET", "content", HandleContentInfoMultiple)
	router.AddMethodHandler("GET", "product", HandleProductInfoListByIds)
	router.AddMethodHandler("GET", "pass", HandlePassInfoListByIds)

	contentRouter := apirouter.NewRouter()
	contentRouter.AddMethodHandlerWildcard("GET", "contentId", HandleContentInfo)
	router.AddRouter("content", contentRouter)

	specificContentRouter := apirouter.NewRouter()
	specificContentRouter.AddMethodHandler("GET", "icon", HandleContentIconDownload)
	contentRouter.AddRouterWildcard("contentId", specificContentRouter)

	productRouter := apirouter.NewRouter()
	productRouter.AddMethodHandler("GET", "list", HandleProductInfoList)
	productRouter.AddMethodHandler("GET", "accesses", HandleAccessProductInfoList, requireAuthMiddleware)
	productRouter.AddMethodHandlerWildcard("GET", "productId", HandleProductInfo)
	router.AddRouter("product", productRouter)

	specificProductRouter := apirouter.NewRouter()
	specificProductRouter.AddMethodHandler("GET", "icon", HandleProductIconDownload)
	specificProductRouter.AddMethodHandler("GET", "access", HandleAccessProductInfo, requireAuthMiddleware)
	productRouter.AddRouterWildcard("productId", specificProductRouter)

	passRouter := apirouter.NewRouter()
	passRouter.AddMethodHandler("GET", "list", HandlePassInfoList)
	passRouter.AddMethodHandler("GET", "accesses", HandleAccessPassInfoList, requireAuthMiddleware)
	passRouter.AddMethodHandlerWildcard("GET", "passId", HandlePassInfo)
	router.AddRouter("pass", passRouter)

	specificPassRouter := apirouter.NewRouter()
	specificPassRouter.AddMethodHandler("GET", "icon", HandlePassIconDownload)
	specificPassRouter.AddMethodHandler("GET", "access", HandleAccessPassInfo, requireAuthMiddleware)
	passRouter.AddRouterWildcard("passId", specificPassRouter)

	return router
}
