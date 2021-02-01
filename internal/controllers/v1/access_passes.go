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

type accessPassInfoListResponseBody struct {
	Passes []*accessPassInfo `json:"passes"`
}

type accessPassInfo struct {
	Access         bool                  `json:"access"`
	PassID         string                `json:"passId"`
	ExpirationDate timeutils.EpochTimeMS `json:"expirationDate,omitempty"`
}

// HandleAccessPassInfoList handles pass access info list requests.
func HandleAccessPassInfoList(c echo.Context) error {
	cc := c.(*authmiddlewares.AuthContext)
	accountID := cc.Session.Data.AccountID

	hub := sentryecho.GetHubFromContext(c)
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID: accountID,
		})
	})
	passAccessVOList, err := globals.PassAccessService.GetPassAccessVOListByAccountID(accountID)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	}
	accessPassItems := make([]*accessPassInfo, len(passAccessVOList))
	for i, passAccessVO := range passAccessVOList {
		accessPassItems[i] = &accessPassInfo{
			Access:         true,
			PassID:         passAccessVO.PassID,
			ExpirationDate: passAccessVO.ExpirationDate,
		}
	}

	return c.JSON(http.StatusOK, &accessPassInfoListResponseBody{
		Passes: accessPassItems,
	})
}

// HandleAccessPassInfo handles pass access info requests.
func HandleAccessPassInfo(c echo.Context) error {
	cc := c.(*authmiddlewares.AuthContext)
	accountID := cc.Session.Data.AccountID

	hub := sentryecho.GetHubFromContext(c)
	hub.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetUser(sentry.User{
			ID: accountID,
		})
	})

	passID := c.Param("passId")
	if len(passID) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}

	passAccessVO, err := globals.PassAccessService.GetPassAccessVOByAccountIDPassID(accountID, passID)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	} else if passAccessVO == nil {
		return c.JSON(http.StatusOK, &accessPassInfo{
			Access: false,
			PassID: passID,
		})
	} else {
		return c.JSON(http.StatusOK, &accessPassInfo{
			Access:         true,
			PassID:         passAccessVO.PassID,
			ExpirationDate: passAccessVO.ExpirationDate,
		})
	}
}
