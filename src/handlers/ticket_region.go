package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-info/serverinfo"
	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
)

type ticketRegionResponseBody struct {
	RegionID string `json:"regionId"`
}

// HandleTicketRegion handles ticket region requests.
func HandleTicketRegion(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	ticketID, _ := req.GetPathParam("ticketId")
	if len(ticketID) < serverinfo.ShortStageSize+serverinfo.ShortRegionSize {
		return resp.SetClientError(apierrors.ErrorInvalidParameters.WithField("ticketId"))
	}

	ticketRegion := ticketID[len(ticketID)-serverinfo.ShortRegionSize : len(ticketID)]

	regionID, err := serverinfo.GetRegionNameOfShortName(ticketRegion)
	if err != nil {
		return resp.SetClientError(apierrors.ErrorInvalidParameters.WithField("ticketId"))
	}

	resp.SetBody(&ticketRegionResponseBody{
		RegionID: regionID,
	})

	return nil
}
