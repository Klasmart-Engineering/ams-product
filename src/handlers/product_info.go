package handlers

import (
	"context"
	"strings"

	"bitbucket.org/calmisland/go-server-product/productid"
	"bitbucket.org/calmisland/go-server-product/productservice"
	"bitbucket.org/calmisland/go-server-reference-data/productdata"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
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
	UpdatedDate timeutils.EpochTimeMS   `json:"updateTm"`
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
	productIDs := req.GetQueryParamMulti("id")
	if len(productIDs) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	// NOTE: This is for backwards compatibility, to support comma-separated values in just one parameter
	if len(productIDs) == 1 && strings.ContainsRune(productIDs[0], ',') {
		productIDs = strings.Split(productIDs[0], ",")
	}

	productVOList, err := globals.ProductService.GetProductVOListByIds(productIDs)
	if err != nil {
		return resp.SetServerError(err)
	}

	products := make([]*productInfoResponseBody, len(productVOList))
	for i, productVO := range productVOList {
		appInfo := convertProductAppInfoFromService(productVO.AppInfo)
		products[i] = &productInfoResponseBody{
			ProductID:   productVO.ProductID,
			Title:       productVO.Title,
			Type:        productVO.Type,
			Description: productVO.Description,
			AppInfo:     appInfo,
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
	productID, _ := req.GetPathParam("productId")
	if len(productID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	productVO, err := globals.ProductService.GetProductVOByProductID(productID)
	if err != nil {
		return resp.SetServerError(err)
	} else if productVO == nil {
		return resp.SetClientError(apierrors.ErrorItemNotFound)
	}

	appInfo := convertProductAppInfoFromService(productVO.AppInfo)
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
	productID, _ := req.GetPathParam("productId")
	if len(productID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	productID, err := productid.GetProductIDShortPrefix(productID)
	if err != nil {
		return resp.SetClientError(apierrors.ErrorInputInvalidFormat)
	}

	fileURL, err := services.GetProductIconURL(productID)
	if err != nil {
		return resp.SetServerError(err)
	}

	resp.Redirect(fileURL)
	return nil
}

func convertProductAppInfoFromService(serviceInfo *productservice.ProductAppInfo) *productAppInfo {
	if serviceInfo == nil {
		return nil
	}

	appStoreInfo := convertProductAppStoreInfoFromService(serviceInfo.AppStore)
	googlePlayInfo := convertProductAppStoreInfoFromService(serviceInfo.GooglePlay)
	return &productAppInfo{
		AppStore:   appStoreInfo,
		GooglePlay: googlePlayInfo,
	}
}

func convertProductAppStoreInfoFromService(serviceInfo *productservice.ProductAppStoreInfo) *productAppStoreInfo {
	if serviceInfo == nil {
		return nil
	}

	return &productAppStoreInfo{
		AppID:    serviceInfo.AppID,
		StoreURL: serviceInfo.StoreURL,
	}
}
