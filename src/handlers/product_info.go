package handlers

import (
	"context"
	"strings"

	"bitbucket.org/calmisland/go-server-product/productservice"
	"bitbucket.org/calmisland/go-server-reference-data/productdata"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/src/services"
)

type productInfoListResponseBody struct {
	Products []*productInfoResponseBody `json:"products"`
}

type productInfoResponseBody struct {
	ProductID   string                  `json:"prodId"`
	Title       string                  `json:"title"`
	Type        productdata.ProductType `json:"type"`
	Description string                  `json:"description"`
	AppInfo     *productAppInfo         `json:"appInfo,omitempty"`
	UpdatedDate timeutils.UnixTime      `json:"updateTm"`
}

type productAppInfo struct {
	AppStore   *productAppStoreInfo `json:"appStore,omitempty"`
	GooglePlay *productAppStoreInfo `json:"googlePlay,omitempty"`
}

type productAppStoreInfo struct {
	AppID    string `json:"appId"`
	StoreURL string `json:"storeUrl"`
}

// HandleProductInfoListByIds handles product information list requests.
func HandleProductInfoListByIds(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	productIDs, hasQueryParam := req.GetQueryParam("id")
	if !hasQueryParam {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	productVOList, err := productservice.ProductService.GetProductVOListByIds(strings.Split(productIDs, ","))
	if err != nil {
		return resp.SetServerError(err)
	} else if productVOList == nil || (productVOList != nil && len(productVOList) == 0) {
		return resp.SetClientError(apierrors.ErrorItemNotFound)
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

	response := productInfoListResponseBody{
		Products: products,
	}
	resp.SetBody(&response)
	return nil
}

// HandleProductInfo handles product information requests.
func HandleProductInfo(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	productID, hasQueryParam := req.GetPathParam("productId")
	if !hasQueryParam {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	productVO, err := productservice.ProductService.GetProductVOByProductID(productID)
	if err != nil {
		return resp.SetServerError(err)
	} else if productVO == nil {
		return resp.SetClientError(apierrors.ErrorItemNotFound)
	}

	var appInfo *productAppInfo
	var appStoreInfo *productAppStoreInfo
	var googlePlayInfo *productAppStoreInfo

	if productVO.AppInfo != nil {
		if productVO.AppInfo.AppStore != nil {
			appStoreInfo = &productAppStoreInfo{
				AppID:    productVO.AppInfo.AppStore.AppID,
				StoreURL: productVO.AppInfo.AppStore.StoreURL,
			}
		}
		if productVO.AppInfo.GooglePlay != nil {
			googlePlayInfo = &productAppStoreInfo{
				AppID:    productVO.AppInfo.AppStore.AppID,
				StoreURL: productVO.AppInfo.AppStore.StoreURL,
			}
		}
		appInfo = &productAppInfo{
			AppStore:   appStoreInfo,
			GooglePlay: googlePlayInfo,
		}
	}

	response := productInfoResponseBody{
		ProductID:   productVO.ProductID,
		Title:       productVO.Title,
		Type:        productVO.Type,
		Description: productVO.Description,
		AppInfo:     appInfo,
		UpdatedDate: productVO.UpdatedDate,
	}
	resp.SetBody(&response)
	return nil
}

// HandleProductIconDownload handles downloading product icons.
func HandleProductIconDownload(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	productID, hasQueryParam := req.GetPathParam("productId")
	if !hasQueryParam {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	fileURL, err := services.GetProgramIconURL(productID)
	if err != nil {
		return resp.SetServerError(err)
	}

	resp.Redirect(fileURL)
	return nil
}
