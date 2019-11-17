package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-product/passaccessservice"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
)

type accessPassIdListResponseBody struct {
	Passes []*accessPassItem `json:"passes"`
}

type accessPassItem struct {
	PassID         string                `json:"passId"`
	ExpirationDate timeutils.EpochTimeMS `json:"expirationDate"`
}

// HandleAccessPassIdList handles pass access list requests.
func HandleAccessPassIdList(_ context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	accountID := req.Session.Data.AccountID
	passAccessVOList, err := passaccessservice.PassAccessService.GetPassAccessVOListByAccountID(accountID)
	if err != nil {
		resp.SetServerError(err)
	}
	accessPassItems := make([]*accessPassItem, len(passAccessVOList))
	for i, passAccessVO := range passAccessVOList {
		accessPassItems[i] = &accessPassItem{
			PassID:         passAccessVO.PassID,
			ExpirationDate: passAccessVO.ExpirationDate,
		}
	}
	resp.SetBody(&accessPassIdListResponseBody{
		Passes: accessPassItems,
	})
	return nil
}
