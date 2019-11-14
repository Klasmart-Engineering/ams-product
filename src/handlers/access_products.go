package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-product/productaccessservice"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
)

type accessProductIdListResponseBody struct {
	Products []*accessProductItem `json:"products"`
}

type accessProductItem struct {
	ProductID      string                `json:"productId"`
	ExpirationDate timeutils.EpochTimeMS `json:"expirationDate"`
}

// HandleAccessProductIdList handles product access list requests.
func HandleAccessProductIdList(_ context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	accountID := req.Session.Data.AccountID
	productAccessVOList, err := productaccessservice.ProductAccessService.GetProductAccessVOListByAccountID(accountID)
	if err != nil {
		resp.SetServerError(err)
	}
	accessProductItems := make([]*accessProductItem, len(productAccessVOList))
	for i, productAccessVO := range productAccessVOList {
		accessProductItems[i] = &accessProductItem{
			ProductID:      productAccessVO.ProductID,
			ExpirationDate: productAccessVO.ExpirationDate,
		}
	}
	resp.SetBody(&accessProductIdListResponseBody{
		Products: accessProductItems,
	})
	return nil
}
