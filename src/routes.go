package main

import (
	"net/http"
	"strings"

	"bitbucket.org/calmisland/go-server-shared/requests/apirequests"
	"bitbucket.org/calmisland/go-server-shared/requests/httprequests"
	"bitbucket.org/calmisland/go-server-shared/services/aws/awslambda"
	"bitbucket.org/calmisland/product-lambda-funcs/src/handlers"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

const (
	requestIDHeaderName = "X-Req-ID"
)

var (
	rootRouter *apirequests.Router
)

func initLambdaFunctions() {
	rootRouter = &apirequests.Router{}
	routerV1 := createLambdaRouterV1()
	rootRouter.AddRouter("v1", routerV1)
}

func createLambdaRouterV1() *apirequests.Router {
	router := &apirequests.Router{}
	router.AddMethodHandler("GET", "product", handlers.HandleProductInfoListByIds)
	router.AddMethodHandler("GET", "serverinfo", handlers.HandleServerInfo)

	productRouter := &apirequests.Router{}
	productRouter.AddMethodHandlerWildcard("GET", "productId", handlers.HandleProductInfo)
	router.AddRouter("product", productRouter)

	return router
}

func handleHTTPRequest(w http.ResponseWriter, r *http.Request) error {
	var apiResponse *apirequests.Response
	urlPath := strings.TrimPrefix(r.URL.Path, "/")
	apiRequest, err := httprequests.ConvertToAPIRequest(r)
	if err == nil {
		apiResponse, err = rootRouter.HandleRequest(apiRequest, urlPath)
		if err != nil {
			apiResponse, err = apirequests.ServerError(err)
		}
	} else {
		apiResponse, err = apirequests.ServerError(err)
	}

	if len(apiRequest.RequestID) > 0 {
		w.Header().Set(requestIDHeaderName, apiRequest.RequestID)
	}
	httprequests.WriteAPIResponse(w, apiResponse)
	return err
}

func handleLambdaRequest(ctx *lambdacontext.LambdaContext, req *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	urlPath := strings.TrimPrefix(req.Path, "/")
	apiRequest := awslambda.ConvertToAPIRequest(ctx, req)
	apiResponse, err := rootRouter.HandleRequest(apiRequest, urlPath)
	if err != nil {
		apiResponse, err = apirequests.ServerError(err)
	}

	resp := awslambda.ConvertToAPIGatewayResponse(apiResponse)

	if len(apiRequest.RequestID) > 0 {
		if resp.MultiValueHeaders != nil {
			resp.MultiValueHeaders[requestIDHeaderName] = []string{
				apiRequest.RequestID,
			}
		} else if resp.Headers != nil {
			resp.Headers[requestIDHeaderName] = apiRequest.RequestID
		} else {
			resp.Headers = map[string]string{
				requestIDHeaderName: apiRequest.RequestID,
			}
		}
	}

	return resp, err
}
