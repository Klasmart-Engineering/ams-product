package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
)

type passInfoListResponseBody struct {
	Passes []*passInfoResponseBody `json:"passes"`
}

type passInfoResponseBody struct {
	PassID     string   `json:"passId"`
	ProductIDs []string `json:"productIds"`
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

	passes := make([]*passInfoResponseBody, len(passVOList))
	for i, passVO := range passVOList {
		passes[i] = &passInfoResponseBody{
			PassID:     passVO.PassID,
			ProductIDs: passVO.Products,
		}
	}

	response := passInfoListResponseBody{
		Passes: passes,
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

	response := passInfoResponseBody{
		PassID:     passVO.PassID,
		ProductIDs: passVO.Products,
	}
	resp.SetBody(&response)
	return nil
}
