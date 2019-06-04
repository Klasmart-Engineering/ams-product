package handlers

import (
	"strings"

	"bitbucket.org/calmisland/go-server-shared/apierrors"
	"bitbucket.org/calmisland/go-server-shared/datareference"
	"bitbucket.org/calmisland/go-server-shared/requests/apirequests"
	"bitbucket.org/calmisland/go-server-shared/utils/timeutils"
	"bitbucket.org/calmisland/go-server-standard/stdservices"
)

type productInfoListResponseBody struct {
	Products []*productInfoResponseBody `json:"products"`
}

type productInfoResponseBody struct {
	ProductID   string                    `json:"prodId"`
	Title       string                    `json:"title"`
	Type        datareference.ProductType `json:"type"`
	Description string                    `json:"description"`
	UpdatedDate timeutils.UnixTime        `json:"updateTm"`
}

// HandleProductInfoListByIds handles product information list requests.
func HandleProductInfoListByIds(req *apirequests.Request) (*apirequests.Response, error) {
	if req.HTTPMethod != "GET" {
		return apirequests.ClientError(apierrors.ErrorBadRequestMethod)
	}

	sessionData := req.ValidateRequestToken()
	if sessionData == nil {
		return apirequests.ClientError(apierrors.ErrorExpiredAccessToken)
	}

	productIDs := req.GetQueryParam("id")
	if productIDs == nil {
		return apirequests.ClientError(apierrors.ErrorInvalidParameters)
	}

	productVOList, err := stdservices.ProductService.GetProductVOListByIds(strings.Split(*productIDs, ","))
	if err != nil {
		return apirequests.ServerError(err)
	} else if productVOList == nil || (productVOList != nil && len(productVOList) == 0) {
		return apirequests.ClientError(apierrors.ErrorItemNotFound)
	}

	products := make([]*productInfoResponseBody, len(productVOList))
	for i, productVO := range productVOList {
		products[i] = &productInfoResponseBody{
			ProductID:   productVO.ProductID,
			Title:       productVO.Title,
			Type:        productVO.Type,
			Description: productVO.Description,
			UpdatedDate: productVO.UpdatedDate,
		}
	}
	return apirequests.NewResponse(&productInfoListResponseBody{
		Products: products,
	})
}

// HandleProductInfo handles product information requests.
func HandleProductInfo(req *apirequests.Request) (*apirequests.Response, error) {
	if req.HTTPMethod != "GET" {
		return apirequests.ClientError(apierrors.ErrorBadRequestMethod)
	}

	sessionData := req.ValidateRequestToken()
	if sessionData == nil {
		return apirequests.ClientError(apierrors.ErrorExpiredAccessToken)
	}

	productID := req.GetPathParam("productId")
	if productID == nil {
		return apirequests.ClientError(apierrors.ErrorInvalidParameters)
	}

	productVO, err := stdservices.ProductService.GetProductVOByProductID(*productID)
	if err != nil {
		return apirequests.ServerError(err)
	} else if productVO == nil {
		return apirequests.ClientError(apierrors.ErrorItemNotFound)
	}

	response := productInfoResponseBody{
		ProductID:   productVO.ProductID,
		Title:       productVO.Title,
		Type:        productVO.Type,
		Description: productVO.Description,
		UpdatedDate: productVO.UpdatedDate,
	}
	return apirequests.NewResponse(response)
}
