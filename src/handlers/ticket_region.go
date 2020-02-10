package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-info/serverinfo"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
)

type ticketRegionResponseBody struct {
	Region string `json:"region"`
	Stage  string `json:"stage"`
}

// HandleTicketRegion handles ticket region requests.
func HandleTicketRegion(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	ticketID, _ := req.GetPathParam("ticketId")
	if len(ticketID) < serverinfo.ShortStageSize+serverinfo.ShortRegionSize {
		return resp.SetClientError(apierrors.ErrorInvalidParameters.WithField("ticketId"))
	}

	ticketStage := ticketID[len(ticketID)-serverinfo.ShortStageSize : len(ticketID)]
	ticketRegion := ticketID[len(ticketID)-serverinfo.ShortStageSize-serverinfo.ShortRegionSize : len(ticketID)-serverinfo.ShortStageSize]

	stage, err := serverinfo.GetStageNameOfShortName(ticketStage)
	if err != nil {
		return resp.SetClientError(apierrors.ErrorInvalidParameters.WithField("ticketId"))
	}
	region, err := serverinfo.GetRegionNameOfShortName(ticketRegion)
	if err != nil {
		return resp.SetClientError(apierrors.ErrorInvalidParameters.WithField("ticketId"))
	}

	resp.SetBody(&ticketRegionResponseBody{
		Region: region,
		Stage:  stage,
	})

	return nil
}
