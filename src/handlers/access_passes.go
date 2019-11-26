package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
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
func HandleAccessPassInfoList(_ context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	accountID := req.Session.Data.AccountID
	passAccessVOList, err := globals.PassAccessService.GetPassAccessVOListByAccountID(accountID)
	if err != nil {
		resp.SetServerError(err)
	}
	accessPassItems := make([]*accessPassInfo, len(passAccessVOList))
	for i, passAccessVO := range passAccessVOList {
		accessPassItems[i] = &accessPassInfo{
			Access:         true,
			PassID:         passAccessVO.PassID,
			ExpirationDate: passAccessVO.ExpirationDate,
		}
	}
	resp.SetBody(&accessPassInfoListResponseBody{
		Passes: accessPassItems,
	})
	return nil
}

// HandleAccessPassInfo handles pass access info requests.
func HandleAccessPassInfo(_ context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	passID, _ := req.GetPathParam("passId")
	if len(passID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}
	accountID := req.Session.Data.AccountID
	passAccessVO, err := globals.PassAccessService.GetPassAccessVOByAccountIDPassID(accountID, passID)
	if err != nil {
		return err
	} else if passAccessVO == nil {
		resp.SetBody(&accessPassInfo{
			Access: false,
			PassID: passID,
		})
	} else {
		resp.SetBody(&accessPassInfo{
			Access:         true,
			PassID:         passAccessVO.PassID,
			ExpirationDate: passAccessVO.ExpirationDate,
		})
	}
	return nil
}
