package v1

import (
	"net/http"
	"strings"

	"bitbucket.org/calmisland/go-server-product/productid"
	"bitbucket.org/calmisland/go-server-product/productservice"
	"bitbucket.org/calmisland/go-server-reference-data/productdata"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/helpers"
	v1Services "bitbucket.org/calmisland/product-lambda-funcs/internal/services/v1"
	"github.com/labstack/echo/v4"
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

// HandleProductInfoList handles product information list requests.
func HandleProductInfoList(c echo.Context) error {
	productVOList, err := globals.ProductService.GetProductVOList()
	if err != nil {
		return helpers.HandleInternalError(c, err)
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
	return c.JSON(http.StatusOK, response)
}

// HandleProductInfoListByIds handles product information list requests.
func HandleProductInfoListByIds(c echo.Context) error {
	productIDs, err := helpers.GetArrayQueryParams(c, "id")
	if len(productIDs) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}

	// NOTE: This is for backwards compatibility, to support comma-separated values in just one parameter
	if len(productIDs) == 1 && strings.ContainsRune(productIDs[0], ',') {
		productIDs = strings.Split(productIDs[0], ",")
	}

	productVOList, err := globals.ProductService.GetProductVOListByIds(productIDs)
	if err != nil {
		return helpers.HandleInternalError(c, err)
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
	return c.JSON(http.StatusOK, response)
}

// HandleProductInfo handles product information requests.
func HandleProductInfo(c echo.Context) error {
	productID := c.Param("productId")
	if len(productID) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}

	productVO, err := globals.ProductService.GetProductVOByProductID(productID)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	} else if productVO == nil {
		return apirequests.EchoSetClientError(c, apierrors.ErrorItemNotFound)
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
	return c.JSON(http.StatusOK, response)
}

// HandleProductIconDownload handles downloading product icons.
func HandleProductIconDownload(c echo.Context) error {
	productID := c.Param("productId")

	if len(productID) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}

	productID, err := productid.GetProductIDShortPrefix(productID)
	if err != nil {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInputInvalidFormat)
	}

	fileURL, err := v1Services.GetProductIconURL(productID)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, fileURL)
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
