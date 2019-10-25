package handlers

import (
	"context"

	"bitbucket.org/calmisland/go-server-requests/apierrors"
	"bitbucket.org/calmisland/go-server-requests/apirequests"
)

type activateTicketRequestBody struct {
	TicketID string `json:"ticketId"`
}

// HandleTicketActivate handles ticket activation requests.
func HandleTicketActivate(ctx context.Context, req *apirequests.Request, resp *apirequests.Response) error {

	// Parse the request body
	var reqBody activateTicketRequestBody
	err := req.UnmarshalBody(&reqBody)
	if err != nil {
		return resp.SetClientError(apierrors.ErrorBadRequestBody)
	}

	ticketID := reqBody.TicketID

	if len(ticketID) == 0 {
		return resp.SetClientError(apierrors.ErrorInvalidParameters)
	}

	return nil
}
