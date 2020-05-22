package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-product/passes"
	"bitbucket.org/calmisland/go-server-product/productid"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/go-server-utils/timeutils"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/src/services"
)

type klpPassInfoListResponseBody struct {
	Passes []*klpPassInfoResponseBody `json:"passes"`
}

type klpPassInfoResponseBody struct {
	PassID       string                `json:"passId"`
	Title        string                `json:"title"`
	Price        string                `json:"price"`
	Currency     passes.Currency       `json:"currency"`
	Duration     passes.DurationDays   `json:"duration"`
	Limits       *klpPassLimits        `json:"limits,omitempty"`
	Accesses     *klpPassAccesses      `json:"accesses,omitempty"`
	Customizable bool                  `json:"customizable"`
	CreatedDate  timeutils.EpochTimeMS `json:"createTm"`
	UpdatedDate  timeutils.EpochTimeMS `json:"updateTm"`
}

type klpPassLimits struct {
	Teacher             *int `json:"teacher"`
	Profile             *int `json:"profile"`
	ClassOnlineDuration *int `json:"classOnlineDuration"`
	ClassOnlineProfile  *int `json:"classOnlineProfile"`
	CloudStorage        *int `json:"cloudStorage"`
}

type klpPassAccesses struct {
	CMS           *bool `json:"CMS"`
	LMS           *bool `json:"LMS"`
	OrgManagement *bool `json:"orgManagement"`
	ProfileAdmin  *bool `json:"profileAdmin"`
	Parents       *bool `json:"parents"`
	FullContent   *bool `json:"fullContent"`
}

func HandleKlpPassInfoList(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	passVOList, err := globals.KlpPassService.GetPassVOList()
	if err != nil {
		return resp.SetServerError(err)
	}

	passes := make([]*klpPassInfoResponseBody, len(passVOList))
	for i, passVO := range passVOList {
		price, err := passVO.Price.ToString(passVO.Currency)
		if err != nil {
			return resp.SetServerError(err)
		}
		var limits *klpPassLimits
		if passVO.Limits != nil {
			limits = &klpPassLimits{
				Teacher:             passVO.Limits.Teacher,
				Profile:             passVO.Limits.Profile,
				ClassOnlineDuration: passVO.Limits.ClassOnlineDuration,
				ClassOnlineProfile:  passVO.Limits.ClassOnlineProfile,
				CloudStorage:        passVO.Limits.CloudStorage,
			}
		}
		var accesses *klpPassAccesses
		if passVO.Accesses != nil {
			accesses = &klpPassAccesses{
				CMS:           passVO.Accesses.CMS,
				LMS:           passVO.Accesses.LMS,
				OrgManagement: passVO.Accesses.OrgManagement,
				ProfileAdmin:  passVO.Accesses.ProfileAdmin,
				Parents:       passVO.Accesses.Parents,
				FullContent:   passVO.Accesses.FullContent,
			}
		}
		passes[i] = &klpPassInfoResponseBody{
			PassID:       passVO.PassID,
			Title:        passVO.Title,
			Price:        price,
			Currency:     passVO.Currency,
			Duration:     passVO.Duration,
			Limits:       limits,
			Accesses:     accesses,
			Customizable: passVO.Customizable,
			CreatedDate:  passVO.CreatedDate,
			UpdatedDate:  passVO.UpdatedDate,
		}
	}

	response := klpPassInfoListResponseBody{
		Passes: passes,
	}
	resp.SetBody(&response)
	return nil
}

// HandleKlpPassInfoListByIds handles kidsloop pass information list requests.
func HandleKlpPassInfoListByIds(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	passIDs := req.GetQueryParamMulti("id")
	if len(passIDs) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	passVOList, err := globals.KlpPassService.GetPassVOListByIds(passIDs)
	if err != nil {
		return resp.SetServerError(err)
	}

	passes := make([]*klpPassInfoResponseBody, len(passVOList))
	for i, passVO := range passVOList {
		price, err := passVO.Price.ToString(passVO.Currency)
		if err != nil {
			return resp.SetServerError(err)
		}
		var limits *klpPassLimits
		if passVO.Limits != nil {
			limits = &klpPassLimits{
				Teacher:             passVO.Limits.Teacher,
				Profile:             passVO.Limits.Profile,
				ClassOnlineDuration: passVO.Limits.ClassOnlineDuration,
				ClassOnlineProfile:  passVO.Limits.ClassOnlineProfile,
				CloudStorage:        passVO.Limits.CloudStorage,
			}
		}
		var accesses *klpPassAccesses
		if passVO.Accesses != nil {
			accesses = &klpPassAccesses{
				CMS:           passVO.Accesses.CMS,
				LMS:           passVO.Accesses.LMS,
				OrgManagement: passVO.Accesses.OrgManagement,
				ProfileAdmin:  passVO.Accesses.ProfileAdmin,
				Parents:       passVO.Accesses.Parents,
				FullContent:   passVO.Accesses.FullContent,
			}
		}
		passes[i] = &klpPassInfoResponseBody{
			PassID:       passVO.PassID,
			Title:        passVO.Title,
			Price:        price,
			Currency:     passVO.Currency,
			Duration:     passVO.Duration,
			Limits:       limits,
			Accesses:     accesses,
			Customizable: passVO.Customizable,
			CreatedDate:  passVO.CreatedDate,
			UpdatedDate:  passVO.UpdatedDate,
		}
	}

	response := klpPassInfoListResponseBody{
		Passes: passes,
	}
	resp.SetBody(&response)
	return nil
}

// HandleKlpPassInfo handles kidsloop pass information requests.
func HandleKlpPassInfo(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	passID, _ := req.GetPathParam("passId")
	if len(passID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters.WithField("passId"))
	}

	passVO, err := globals.KlpPassService.GetPassVOByPassID(passID)
	if err != nil {
		return resp.SetServerError(err)
	} else if passVO == nil {
		return resp.SetClientError(apierrors.ErrorItemNotFound)
	}
	price, err := passVO.Price.ToString(passVO.Currency)
	if err != nil {
		return resp.SetServerError(err)
	}
	var limits *klpPassLimits
	if passVO.Limits != nil {
		limits = &klpPassLimits{
			Teacher:             passVO.Limits.Teacher,
			Profile:             passVO.Limits.Profile,
			ClassOnlineDuration: passVO.Limits.ClassOnlineDuration,
			ClassOnlineProfile:  passVO.Limits.ClassOnlineProfile,
			CloudStorage:        passVO.Limits.CloudStorage,
		}
	}
	var accesses *klpPassAccesses
	if passVO.Accesses != nil {
		accesses = &klpPassAccesses{
			CMS:           passVO.Accesses.CMS,
			LMS:           passVO.Accesses.LMS,
			OrgManagement: passVO.Accesses.OrgManagement,
			ProfileAdmin:  passVO.Accesses.ProfileAdmin,
			Parents:       passVO.Accesses.Parents,
			FullContent:   passVO.Accesses.FullContent,
		}
	}
	response := &klpPassInfoResponseBody{
		PassID:       passVO.PassID,
		Title:        passVO.Title,
		Price:        price,
		Currency:     passVO.Currency,
		Duration:     passVO.Duration,
		Limits:       limits,
		Accesses:     accesses,
		Customizable: passVO.Customizable,
		CreatedDate:  passVO.CreatedDate,
		UpdatedDate:  passVO.UpdatedDate,
	}
	resp.SetBody(&response)
	return nil
}

// HandleKlpPassIconDownload handles downloading kidsloop pass icons.
func HandleKlpPassIconDownload(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	passID, _ := req.GetPathParam("passId")
	if len(passID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	passID, err := productid.GetPassIDShortPrefix(passID)
	if err != nil {
		return resp.SetClientError(apierrors.ErrorInputInvalidFormat)
	}

	fileURL, err := services.GetKlpPassIconURL(passID)
	if err != nil {
		return resp.SetServerError(err)
	}

	resp.Redirect(fileURL)
	return nil
}
