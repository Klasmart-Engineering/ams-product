package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-product/passes"
	"bitbucket.org/calmisland/go-server-product/productid"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/src/services"
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
func HandlePassInfoList(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	passVOList, err := globals.PassService.GetPassVOList()

	if err != nil {
		return resp.SetServerError(err)
	}

	passesObj := make([]*passInfoResponseBody, len(passVOList))
	for i, passVO := range passVOList {
		price, err := passVO.Price.ToString(passVO.Currency)
		if err != nil {
			return resp.SetServerError(err)
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
	resp.SetBody(&response)
	return nil
}

// HandlePassInfoListByIds handles pass information list requests.
func HandlePassInfoListByIds(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	passIDs := req.GetQueryParamMulti("id")
	if len(passIDs) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	passVOList, err := globals.PassService.GetPassVOListByIds(passIDs)
	if err != nil {
		return resp.SetServerError(err)
	}

	passesObj := make([]*passInfoResponseBody, len(passVOList))
	for i, passVO := range passVOList {
		price, err := passVO.Price.ToString(passVO.Currency)
		if err != nil {
			return resp.SetServerError(err)
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
	resp.SetBody(&response)
	return nil
}

// HandlePassInfo handles pass information requests.
func HandlePassInfo(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	passID, _ := req.GetPathParam("passId")
	if len(passID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters.WithField("passId"))
	}

	passVO, err := globals.PassService.GetPassVOByPassID(passID)
	if err != nil {
		return resp.SetServerError(err)
	} else if passVO == nil {
		return resp.SetClientError(apierrors.ErrorItemNotFound)
	}
	price, err := passVO.Price.ToString(passVO.Currency)
	if err != nil {
		return resp.SetServerError(err)
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

	resp.SetBody(&response)
	return nil
}

// HandlePassIconDownload handles downloading pass icons.
func HandlePassIconDownload(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	passID, _ := req.GetPathParam("passId")
	if len(passID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	passID, err := productid.GetPassIDShortPrefix(passID)
	if err != nil {
		return resp.SetClientError(apierrors.ErrorInputInvalidFormat)
	}

	fileURL, err := services.GetPassIconURL(passID)
	if err != nil {
		return resp.SetServerError(err)
	}

	resp.Redirect(fileURL)
	return nil
}
