package handlers

import (
	"bitbucket.org/calmisland/go-server-shared/apierrors"
	"bitbucket.org/calmisland/go-server-shared/requests/apirequests"
	"bitbucket.org/calmisland/go-server-shared/serverinfo"
)

type serverInfoResponseBody struct {
	RegionName string `json:"region"`
	StageName  string `json:"stage"`
}

// HandleServerInfo handles server information requests.
func HandleServerInfo(req *apirequests.Request) (*apirequests.Response, error) {
	if req.HTTPMethod != "GET" {
		return apirequests.ClientError(apierrors.ErrorBadRequestMethod)
	}

	serverRegionName := serverinfo.GetRegionName()
	serverStageName := serverinfo.GetStageName()
	response := serverInfoResponseBody{
		RegionName: serverRegionName,
		StageName:  serverStageName,
	}
	return apirequests.NewResponse(response)
}
