package handlers

import (
	"bitbucket.org/calmisland/go-server-shared/v2/apierrors"
	"bitbucket.org/calmisland/go-server-shared/v2/requests/apirequests"
	"bitbucket.org/calmisland/go-server-shared/v2/serverinfo"
)

type serverInfoResponseBody struct {
	RegionName string `json:"region"`
	StageName  string `json:"stage"`
}

// HandleServerInfo handles server information requests.
func HandleServerInfo(req *apirequests.Request) (*apirequests.Response, error) {
	if req.HTTPMethod != "GET" {
		return apirequests.ClientError(req, apierrors.ErrorBadRequestMethod)
	}

	serverRegionName := serverinfo.GetRegionName()
	serverStageName := serverinfo.GetStageName()
	response := serverInfoResponseBody{
		RegionName: serverRegionName,
		StageName:  serverStageName,
	}
	return apirequests.NewResponse(response)
}
