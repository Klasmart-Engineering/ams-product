package v1

import (
	"net/http"

	"bitbucket.org/calmisland/go-server-auth/authmiddlewares"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/helpers"
	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
)

type accessProductInfoListResponseBody struct {
	Products []*accessProductInfo `json:"products"`
}

type accessProductInfo struct {
	Access         bool                  `json:"access"`
	ProductID      string                `json:"productId"`
	ExpirationDate timeutils.EpochTimeMS `json:"expirationDate,omitempty"`
}

// HandleAccessProductInfoList handles product access info list requests.
func HandleAccessProductInfoList(c echo.Context) error {
	cc := c.(*authmiddlewares.AuthContext)
	accountID := cc.Session.Data.AccountID

	hub := sentryecho.GetHubFromContext(c)
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID: accountID,
		})
	})

	productAccessVOList, err := globals.ProductAccessService.GetProductAccessVOListByAccountID(accountID)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	}
	accessProductItems := make([]*accessProductInfo, len(productAccessVOList))
	for i, productAccessVO := range productAccessVOList {
		accessProductItems[i] = &accessProductInfo{
			Access:         true,
			ProductID:      productAccessVO.ProductID,
			ExpirationDate: productAccessVO.ExpirationDate,
		}
	}

	return c.JSON(http.StatusOK, &accessProductInfoListResponseBody{
		Products: accessProductItems,
	})
}

// HandleAccessProductInfo handles product access info requests.
func HandleAccessProductInfo(c echo.Context) error {
	cc := c.(*authmiddlewares.AuthContext)
	accountID := cc.Session.Data.AccountID

	hub := sentryecho.GetHubFromContext(c)
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID: accountID,
		})
	})

	productID := c.Param("productId")
	if len(productID) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}
	productAccessVO, err := globals.ProductAccessService.GetProductAccessVOByAccountIDProductID(accountID, productID)
	if err != nil {
		return err
	} else if productAccessVO == nil {
		return c.JSON(http.StatusOK, &accessProductInfo{
			Access:    false,
			ProductID: productID,
		})
	} else {
		return c.JSON(http.StatusOK, &accessProductInfo{
			Access:         true,
			ProductID:      productAccessVO.ProductID,
			ExpirationDate: productAccessVO.ExpirationDate,
		})
	}
}
