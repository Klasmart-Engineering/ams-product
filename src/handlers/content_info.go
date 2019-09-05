package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/src/services"
)

type contentInfoListResponseBody struct {
	Contents []*contentInfoResponseBody `json:"contents"`
}

type contentInfoResponseBody struct {
	ContentID   string              `json:"id"`
	ProductID   string              `json:"productId"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	KidsAppInfo *kidsAppContentInfo `json:"kidsAppInfo,omitempty"`
	UpdatedDate timeutils.UnixTime  `json:"updateTm"`
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

	response := contentInfoResponseBody{}
	resp.SetBody(&response)
	return nil
}

// HandleContentInfoMultiple handles multiple content informations requests.
func HandleContentInfoMultiple(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	contentIDs := req.GetQueryParamMulti("id")
	if len(contentIDs) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	response := contentInfoListResponseBody{}
	resp.SetBody(&response)
	return nil
}

// HandleContentIconDownload handles content icon download requests.
func HandleContentIconDownload(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	contentID, _ := req.GetPathParam("contentId")
	if len(contentID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	fileURL, err := services.GetContentIconURL(contentID)
	if err != nil {
		return resp.SetServerError(err)
	}

	resp.Redirect(fileURL)
	return nil
}
