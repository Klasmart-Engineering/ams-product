package handlers

import (
	"context"

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
	resp.SetBody(&accessProductIdListResponseBody{
		Products: []*accessProductItem{
			&accessProductItem{
				ProductID:      "app.learnandplay.bada-genius",
				ExpirationDate: timeutils.EpochMSNow(),
			},
		},
	})
	return nil
}
