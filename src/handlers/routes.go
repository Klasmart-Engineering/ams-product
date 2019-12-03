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
	routerV1 := createLambdaRouterV1()
	rootRouter.AddRouter("v1", routerV1)
	return rootRouter
}

func createLambdaRouterV1() *apirouter.Router {
	router := apirouter.NewRouter()
	router.AddMethodHandler("GET", "serverinfo", standardhandlers.HandleServerInfo)

	requireAuthMiddleware := authmiddlewares.ValidateSession(globals.AccessTokenValidator, true)

	router.AddMethodHandler("GET", "content", HandleContentInfoMultiple, requireAuthMiddleware)
	router.AddMethodHandler("GET", "product", HandleProductInfoListByIds, requireAuthMiddleware)
	router.AddMethodHandler("GET", "pass", HandlePassInfoListByIds, requireAuthMiddleware)

	contentRouter := apirouter.NewRouter()
	contentRouter.AddMiddleware(requireAuthMiddleware)
	contentRouter.AddMethodHandlerWildcard("GET", "contentId", HandleContentInfo)
	router.AddRouter("content", contentRouter)

	specificContentRouter := apirouter.NewRouter()
	specificContentRouter.AddMethodHandler("GET", "icon", HandleContentIconDownload)
	contentRouter.AddRouterWildcard("contentId", specificContentRouter)

	productRouter := apirouter.NewRouter()
	productRouter.AddMiddleware(requireAuthMiddleware)
	productRouter.AddMethodHandler("GET", "accesses", HandleAccessProductInfoList)
	productRouter.AddMethodHandlerWildcard("GET", "productId", HandleProductInfo)
	router.AddRouter("product", productRouter)

	specificProductRouter := apirouter.NewRouter()
	specificProductRouter.AddMethodHandler("GET", "icon", HandleProductIconDownload)
	specificProductRouter.AddMethodHandler("GET", "access", HandleAccessProductInfo)
	productRouter.AddRouterWildcard("productId", specificProductRouter)

	passRouter := apirouter.NewRouter()
	passRouter.AddMiddleware(requireAuthMiddleware)
	passRouter.AddMethodHandler("GET", "accesses", HandleAccessPassInfoList)
	passRouter.AddMethodHandlerWildcard("GET", "passId", HandlePassInfo)
	router.AddRouter("pass", passRouter)

	specificPassRouter := apirouter.NewRouter()
	specificPassRouter.AddMethodHandler("GET", "access", HandleAccessPassInfo)
	passRouter.AddRouterWildcard("passId", specificPassRouter)

	return router
}
