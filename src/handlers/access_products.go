package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
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
func HandleAccessProductInfoList(_ context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	accountID := req.Session.Data.AccountID
	productAccessVOList, err := globals.ProductAccessService.GetProductAccessVOListByAccountID(accountID)
	if err != nil {
		return resp.SetServerError(err)
	}
	accessProductItems := make([]*accessProductInfo, len(productAccessVOList))
	for i, productAccessVO := range productAccessVOList {
		accessProductItems[i] = &accessProductInfo{
			Access:         true,
			ProductID:      productAccessVO.ProductID,
			ExpirationDate: productAccessVO.ExpirationDate,
		}
	}
	resp.SetBody(&accessProductInfoListResponseBody{
		Products: accessProductItems,
	})
	return nil
}

// HandleAccessProductInfo handles product access info requests.
func HandleAccessProductInfo(_ context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	productID, _ := req.GetPathParam("productId")
	if len(productID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}
	accountID := req.Session.Data.AccountID
	productAccessVO, err := globals.ProductAccessService.GetProductAccessVOByAccountIDProductID(accountID, productID)
	if err != nil {
		return err
	} else if productAccessVO == nil {
		resp.SetBody(&accessProductInfo{
			Access:    false,
			ProductID: productID,
		})
	} else {
		resp.SetBody(&accessProductInfo{
			Access:         true,
			ProductID:      productAccessVO.ProductID,
			ExpirationDate: productAccessVO.ExpirationDate,
		})
	}
	return nil
}
