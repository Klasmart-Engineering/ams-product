package globals

import (
	"errors"

	"bitbucket.org/calmisland/go-server-product/contentservice"
	"bitbucket.org/calmisland/go-server-product/klppassservice"
	"bitbucket.org/calmisland/go-server-product/passaccessservice"
	"bitbucket.org/calmisland/go-server-product/passservice"
	"bitbucket.org/calmisland/go-server-product/productaccessservice"
	"bitbucket.org/calmisland/go-server-product/productdatabase"
	"bitbucket.org/calmisland/go-server-product/productservice"
	"bitbucket.org/calmisland/go-server-requests/tokens/accesstokens"
)

var (
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

	// ContentService ContentService
	ContentService contentservice.IContentService

	// KlpPassService KlpPassService
	KlpPassService klppassservice.IKlpPassService
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
	} else if KlpPassService == nil {
		panic(errors.New("The kidsloop pass service has not been set"))
	}
}
