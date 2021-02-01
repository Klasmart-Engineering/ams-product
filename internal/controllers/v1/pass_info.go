package v1

import (
	"net/http"

	"bitbucket.org/calmisland/go-server-product/passes"
	"bitbucket.org/calmisland/go-server-product/productid"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/helpers"
	v1Services "bitbucket.org/calmisland/product-lambda-funcs/internal/services/v1"
	"github.com/labstack/echo/v4"
)

type passInfoListResponseBody struct {
	Passes []*passInfoResponseBody `json:"passes"`
}

type passInfoResponseBody struct {
	PassID     string                      `json:"passId"`
	Title      string                      `json:"title"`
	ProductIDs []string                    `json:"productIds"`
	Price      string                      `json:"price"`
	Currency   passes.Currency             `json:"currency"`
	Duration   passes.DurationDays         `json:"duration"`
	DurationMS passes.DurationMilliseconds `json:"durationMS"`
}

// HandlePassInfoList is the function has /v1/pass/list
func HandlePassInfoList(c echo.Context) error {
	passVOList, err := globals.PassService.GetPassVOList()

	if err != nil {
		return helpers.HandleInternalError(c, err)
	}

	passesObj := make([]*passInfoResponseBody, len(passVOList))
	for i, passVO := range passVOList {
		price, err := passVO.Price.ToString(passVO.Currency)
		if err != nil {
			return helpers.HandleInternalError(c, err)
		}

		durationMS := passVO.DurationMS
		duration := passVO.Duration

		if durationMS == 0 {
			durationMS = passes.DurationMilliseconds(duration) * 86400000
		} else if duration == 0 {
			duration = passes.DurationDays(durationMS / 86400000)
		}

		passesObj[i] = &passInfoResponseBody{
			PassID:     passVO.PassID,
			Title:      passVO.Title,
			ProductIDs: passVO.Products,
			Price:      price,
			Currency:   passVO.Currency,
			Duration:   duration,
			DurationMS: durationMS,
		}
	}

	response := passInfoListResponseBody{
		Passes: passesObj,
	}
	return c.JSON(http.StatusOK, response)
}

// HandlePassInfoListByIds handles pass information list requests.
func HandlePassInfoListByIds(c echo.Context) error {
	passIDs, err := helpers.GetArrayQueryParams(c, "id")

	if len(passIDs) == 0 || err != nil {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}

	passVOList, err := globals.PassService.GetPassVOListByIds(passIDs)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	}

	passesObj := make([]*passInfoResponseBody, len(passVOList))
	for i, passVO := range passVOList {
		price, err := passVO.Price.ToString(passVO.Currency)
		if err != nil {
			return helpers.HandleInternalError(c, err)
		}

		durationMS := passVO.DurationMS
		duration := passVO.Duration

		if durationMS == 0 {
			durationMS = passes.DurationMilliseconds(duration) * 86400000
		} else if duration == 0 {
			duration = passes.DurationDays(durationMS / 86400000)
		}

		passesObj[i] = &passInfoResponseBody{
			PassID:     passVO.PassID,
			Title:      passVO.Title,
			ProductIDs: passVO.Products,
			Price:      price,
			Currency:   passVO.Currency,
			Duration:   duration,
			DurationMS: durationMS,
		}
	}

	response := passInfoListResponseBody{
		Passes: passesObj,
	}
	return c.JSON(http.StatusOK, response)
}

// HandlePassInfo handles pass information requests.
func HandlePassInfo(c echo.Context) error {
	passID := c.Param("passId")
	if len(passID) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters.WithField("passId"))
	}

	passVO, err := globals.PassService.GetPassVOByPassID(passID)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	} else if passVO == nil {
		return apirequests.EchoSetClientError(c, apierrors.ErrorItemNotFound)
	}
	price, err := passVO.Price.ToString(passVO.Currency)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	}

	durationMS := passVO.DurationMS
	duration := passVO.Duration

	if durationMS == 0 {
		durationMS = passes.DurationMilliseconds(duration) * 86400000
	} else if duration == 0 {
		duration = passes.DurationDays(durationMS / 86400000)
	}

	response := passInfoResponseBody{
		PassID:     passVO.PassID,
		Title:      passVO.Title,
		ProductIDs: passVO.Products,
		Price:      price,
		Currency:   passVO.Currency,
		Duration:   duration,
		DurationMS: durationMS,
	}

	return c.JSON(http.StatusOK, response)
}

// HandlePassIconDownload handles downloading pass icons.
func HandlePassIconDownload(c echo.Context) error {
	passID := c.Param("passId")
	if len(passID) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}

	passID, err := productid.GetPassIDShortPrefix(passID)
	if err != nil {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInputInvalidFormat)
	}

	fileURL, err := v1Services.GetPassIconURL(passID)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, fileURL)
}
