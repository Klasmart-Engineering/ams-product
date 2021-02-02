package routers

import (
	"bitbucket.org/calmisland/go-server-auth/authmiddlewares"
	v1Controller "bitbucket.org/calmisland/product-lambda-funcs/internal/controllers/v1"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/globals"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// SetupRouter is ...
func SetupRouter() *echo.Echo {
	// Echo instance
	e := echo.New()

	authMiddleware := authmiddlewares.EchoAuthMiddleware(globals.AccessTokenValidator, true)

	// Middleware
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(sentryecho.New(sentryecho.Options{}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if hub := sentryecho.GetHubFromContext(ctx); hub != nil {
				hub.Scope().SetUser(sentry.User{
					IPAddress: ctx.RealIP(),
				})
			}
			return next(ctx)
		}
	})

	v1 := e.Group("/v1")

	v1.GET("/serverinfo", v1Controller.HandleServerInfo)

	v1.GET("/content", v1Controller.HandleContentInfoMultiple)
	v1.GET("/content/:contentId", v1Controller.HandleContentInfo)
	v1.GET("/content/:contentId/icon", v1Controller.HandleContentIconDownload)

	v1.GET("/product", v1Controller.HandleProductInfoListByIds)
	v1.GET("/product/list", v1Controller.HandleProductInfoList)
	v1.GET("/product/accesses", v1Controller.HandleAccessProductInfoList, authMiddleware)

	v1.GET("/product/:productId", v1Controller.HandleProductInfo)
	v1.GET("/product/:productId/icon", v1Controller.HandleProductIconDownload)
	v1.GET("/product/:productId/access", v1Controller.HandleAccessProductInfo, authMiddleware)

	v1.GET("/pass", v1Controller.HandlePassInfoListByIds)
	v1.GET("/pass/list", v1Controller.HandlePassInfoList)
	v1.GET("/pass/accesses", v1Controller.HandleAccessPassInfoList, authMiddleware)

	v1.GET("/pass/:passId", v1Controller.HandlePassInfo)
	v1.GET("/pass/:passId/icon", v1Controller.HandlePassIconDownload)
	v1.GET("/pass/:passId/access", v1Controller.HandleAccessPassInfo, authMiddleware)

	return e
}
