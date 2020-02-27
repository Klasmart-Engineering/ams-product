package globals

import (
	"errors"

	"bitbucket.org/calmisland/go-server-product/contentservice"
	"bitbucket.org/calmisland/go-server-product/eventticketregionservice"
	"bitbucket.org/calmisland/go-server-product/passaccessservice"
	"bitbucket.org/calmisland/go-server-product/passservice"
	"bitbucket.org/calmisland/go-server-product/productaccessservice"
	"bitbucket.org/calmisland/go-server-product/productdatabase"
	"bitbucket.org/calmisland/go-server-product/productservice"
	"bitbucket.org/calmisland/go-server-requests/apirouter"
	"bitbucket.org/calmisland/go-server-requests/tokens/accesstokens"
)

var (
	// CORSOptions are the CORS options to use for the API.
	CORSOptions *apirouter.CORSOptions

	// AccessTokenValidator is the access token validator.
	AccessTokenValidator accesstokens.Validator

	// ProductDatabase ProductDatabase
	ProductDatabase productdatabase.Database

	// ProductService ProductService
	ProductService productservice.IProductService

	// ProductAccessService ProductAccessService
	ProductAccessService productaccessservice.IProductAccessService

	// PassService PassService
	PassService passservice.IPassService

	// PassAccessService PassAccessService
	PassAccessService passaccessservice.IPassAccessService

	// EventTicketRegionService EventTicketRegionService
	EventTicketRegionService eventticketregionservice.IService

	// ContentService ContentService
	ContentService contentservice.IContentService
)

// Verify verifies if all variables have been properly set.
func Verify() {
	if AccessTokenValidator == nil {
		panic(errors.New("The access token validator has not been set"))
	} else if ProductDatabase == nil {
		panic(errors.New("The product database has not been set"))
	} else if ProductService == nil {
		panic(errors.New("The product service has not been set"))
	} else if ProductAccessService == nil {
		panic(errors.New("The product access service has not been set"))
	} else if PassService == nil {
		panic(errors.New("The pass service has not been set"))
	} else if PassAccessService == nil {
		panic(errors.New("The pass access service has not been set"))
	} else if ContentService == nil {
		panic(errors.New("The content service has not been set"))
	} else if EventTicketRegionService == nil {
		panic(errors.New("The event ticket region service has not been set"))
	} else if CORSOptions == nil {
		panic(errors.New("The CORS definition has not been set"))
	}
}
