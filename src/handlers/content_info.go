package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-product/contentservice"
	"bitbucket.org/calmisland/go-server-product/productid"
	"bitbucket.org/calmisland/go-server-reference-data/productdata"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/src/services"
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
func HandleContentInfo(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	contentID, _ := req.GetPathParam("contentId")
	if len(contentID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	contentVO, err := globals.ContentService.GetContentVOByContentID(contentID)
	if err != nil {
		return resp.SetServerError(err)
	} else if contentVO == nil {
		return resp.SetClientError(apierrors.ErrorItemNotFound)
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
	resp.SetBody(&response)
	return nil
}

// HandleContentInfoMultiple handles multiple content informations requests.
func HandleContentInfoMultiple(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	contentIDs := req.GetQueryParamMulti("id")
	if len(contentIDs) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	contentVOList, err := globals.ContentService.GetContentVOListByIds(contentIDs)
	if err != nil {
		return resp.SetServerError(err)
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
	resp.SetBody(&response)
	return nil
}

// HandleContentIconDownload handles content icon download requests.
func HandleContentIconDownload(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	contentID, _ := req.GetPathParam("contentId")
	if len(contentID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	contentID, err := productid.GetProductIDShortPrefix(contentID)
	if err != nil {
		return resp.SetClientError(apierrors.ErrorInputInvalidFormat)
	}

	fileURL, err := services.GetContentIconURL(contentID)
	if err != nil {
		return resp.SetServerError(err)
	}

	resp.Redirect(fileURL)
	return nil
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
