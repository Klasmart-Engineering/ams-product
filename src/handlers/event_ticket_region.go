package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
)

type eventTicketRegionResponseBody struct {
	RegionID string `json:"regionId"`
}

// HandleEventTicketRegion handles event ticket region requests.
func HandleEventTicketRegion(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {
	eventID, _ := req.GetPathParam("eventId")
	if len(eventID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters.WithField("eventId"))
	}

	eventTicketRegionInfo, err := globals.EventTicketRegionService.GetEventTicketRegion(eventID)
	if err != nil {
		return resp.SetServerError(err)
	} else if eventTicketRegionInfo == nil {
		return resp.SetClientError(apierrors.ErrorItemNotFound.WithField("eventId"))
	}

	resp.SetBody(&ticketRegionResponseBody{
		RegionID: eventTicketRegionInfo.RegionID,
	})

	return nil
}
