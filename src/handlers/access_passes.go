package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-product/passaccessservice"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
)

type accessPassInfoListResponseBody struct {
	Passes   []*accessPassInfo    `json:"passes"`
	Products []*accessProductInfo `json:"products"`
}

type accessPassInfo struct {
	Access         bool                  `json:"access"`
	PassID         string                `json:"passId"`
	Products       []*accessProductInfo  `json:"products,omitempty"`
	ExpirationDate timeutils.EpochTimeMS `json:"expirationDate,omitempty"`
}

// HandleAccessPassInfoList handles pass access info list requests.
func HandleAccessPassInfoList(_ context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	accountID := req.Session.Data.AccountID
	passAccessVOList, err := globals.PassAccessService.GetPassAccessVOListByAccountID(accountID)
	if err != nil {
		return resp.SetServerError(err)
	}
	accessPassItems := make([]*accessPassInfo, len(passAccessVOList))
	for i, passAccessVO := range passAccessVOList {
		accessPassItems[i] = &accessPassInfo{
			Access:         true,
			PassID:         passAccessVO.PassID,
			ExpirationDate: passAccessVO.ExpirationDate,
		}
	}

	productAccessMap, err := computeProductAccessMatrixByPassList(passAccessVOList)
	if err != nil {
		return resp.SetServerError(err)
	}
	accessProductItems := make([]*accessProductInfo, 0, len(productAccessMap))
	for productID, expirationTm := range productAccessMap {
		accessProductItems = append(accessProductItems, &accessProductInfo{
			Access:         true,
			ProductID:      productID,
			ExpirationDate: expirationTm,
		})
	}

	resp.SetBody(&accessPassInfoListResponseBody{
		Passes:   accessPassItems,
		Products: accessProductItems,
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
		return resp.SetServerError(err)
	} else if passAccessVO == nil {
		resp.SetBody(&accessPassInfo{
			Access: false,
			PassID: passID,
		})
	} else {
		productAccessMap, err := computeProductAccessMatrixByPassList([]*passaccessservice.PassAccessVO{passAccessVO})
		if err != nil {
			return resp.SetServerError(err)
		}
		accessProductItems := make([]*accessProductInfo, 0, len(productAccessMap))
		for productID, expirationTm := range productAccessMap {
			accessProductItems = append(accessProductItems, &accessProductInfo{
				Access:         true,
				ProductID:      productID,
				ExpirationDate: expirationTm,
			})
		}
		resp.SetBody(&accessPassInfo{
			Access:         true,
			PassID:         passAccessVO.PassID,
			Products:       accessProductItems,
			ExpirationDate: passAccessVO.ExpirationDate,
		})
	}
	return nil
}

func computeProductAccessMatrixByPassList(passAccessVOList []*passaccessservice.PassAccessVO) (map[string]timeutils.EpochTimeMS, error) {
	passIds := make([]string, len(passAccessVOList))
	for i, passAccessVO := range passAccessVOList {
		passIds[i] = passAccessVO.PassID
	}
	passVOList, err := globals.PassService.GetPassVOListByIds(passIds)
	if err != nil {
		return nil, err
	}
	productExpirationMap := map[string]timeutils.EpochTimeMS{}
	timeNow := timeutils.EpochMSNow()
	for _, passVO := range passVOList {
		for _, productID := range passVO.Products {
			if _, ok := productExpirationMap[productID]; !ok {
				productExpirationMap[productID] = timeNow
			}
			newExpirationDate := timeutils.ConvEpochTimeMS(productExpirationMap[productID].Time().AddDate(0, 0, int(passVO.Duration)))
			if newExpirationDate > productExpirationMap[productID] {
				productExpirationMap[productID] = newExpirationDate
			}
		}
	}
	return productExpirationMap, nil
}
