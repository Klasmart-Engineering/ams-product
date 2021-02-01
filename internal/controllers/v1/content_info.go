package v1

import (
	"fmt"
	"net/http"

	"bitbucket.org/calmisland/go-server-product/contentservice"
	"bitbucket.org/calmisland/go-server-product/productid"
	"bitbucket.org/calmisland/go-server-reference-data/productdata"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/helpers"
	v1Services "bitbucket.org/calmisland/product-lambda-funcs/internal/services/v1"
	"github.com/labstack/echo/v4"
)

type contentInfoListResponseBody struct {
	Contents []*contentInfoResponseBody `json:"contents"`
}

type contentInfoResponseBody struct {
	ContentID   string                  `json:"id"`
	ProductID   string                  `json:"productId"`
	Title       string                  `json:"title"`
	Type        productdata.ContentType `json:"type"`
	Description string                  `json:"description,omitempty"`
	KidsAppInfo *kidsAppContentInfo     `json:"kidsAppInfo,omitempty"`
	UpdatedDate timeutils.EpochTimeMS   `json:"updateTm"`
}

type kidsAppContentInfo struct {
	ContentID   string `json:"contentId"`
	ContentType string `json:"contentType"`
}

// HandleContentInfo handles content information requests.
func HandleContentInfo(c echo.Context) error {
	contentID := c.Param("contentId")
	fmt.Println(contentID)
	if len(contentID) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}

	contentVO, err := globals.ContentService.GetContentVOByContentID(contentID)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	} else if contentVO == nil {
		return apirequests.EchoSetClientError(c, apierrors.ErrorItemNotFound)
	}

	kidsAppInfo := convertContentKidsAppInfoFromService(contentVO.KidsAppInfo)
	response := contentInfoResponseBody{
		ContentID:   contentVO.ContentID,
		ProductID:   contentVO.ProductID,
		Title:       contentVO.Title,
		Type:        contentVO.Type,
		Description: contentVO.Description,
		KidsAppInfo: kidsAppInfo,
		UpdatedDate: contentVO.UpdatedDate,
	}

	return c.JSON(http.StatusOK, response)
}

// HandleContentInfoMultiple handles multiple content informations requests.
func HandleContentInfoMultiple(c echo.Context) error {
	contentIDs, err := helpers.GetArrayQueryParams(c, "id")

	if err != nil || len(contentIDs) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}

	contentVOList, err := globals.ContentService.GetContentVOListByIds(contentIDs)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	}

	contents := make([]*contentInfoResponseBody, len(contentVOList))
	for i, contentVO := range contentVOList {
		kidsAppInfo := convertContentKidsAppInfoFromService(contentVO.KidsAppInfo)
		contents[i] = &contentInfoResponseBody{
			ContentID:   contentVO.ContentID,
			ProductID:   contentVO.ProductID,
			Title:       contentVO.Title,
			Type:        contentVO.Type,
			Description: contentVO.Description,
			KidsAppInfo: kidsAppInfo,
			UpdatedDate: contentVO.UpdatedDate,
		}
	}

	response := contentInfoListResponseBody{
		Contents: contents,
	}
	return c.JSON(http.StatusOK, response)
}

// HandleContentIconDownload handles content icon download requests.
func HandleContentIconDownload(c echo.Context) error {
	contentID := c.Param("contentId")

	if len(contentID) == 0 {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInvalidParameters)
	}

	contentID, err := productid.GetProductIDShortPrefix(contentID)
	if err != nil {
		return apirequests.EchoSetClientError(c, apierrors.ErrorInputInvalidFormat)
	}

	fileURL, err := v1Services.GetContentIconURL(contentID)
	if err != nil {
		return helpers.HandleInternalError(c, err)
	}

	return c.Redirect(http.StatusTemporaryRedirect, fileURL)
}

func convertContentKidsAppInfoFromService(serviceInfo *contentservice.ContentKidsAppInfo) *kidsAppContentInfo {
	if serviceInfo == nil {
		return nil
	}

	return &kidsAppContentInfo{
		ContentID:   serviceInfo.ContentID,
		ContentType: serviceInfo.ContentType,
	}
}
