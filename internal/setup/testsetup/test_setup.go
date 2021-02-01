package testsetup

import (
	"bitbucket.org/calmisland/go-server-product/contentservice"
	"bitbucket.org/calmisland/go-server-product/passaccessservice"
	"bitbucket.org/calmisland/go-server-product/passservice"
	"bitbucket.org/calmisland/go-server-product/productaccessservice"
	"bitbucket.org/calmisland/go-server-product/productdatabase/productmemorydb"
	"bitbucket.org/calmisland/go-server-product/productservice"
	"bitbucket.org/calmisland/go-server-requests/apirouter"
	"bitbucket.org/calmisland/go-server-requests/sessions"
	"bitbucket.org/calmisland/go-server-requests/tokens/accesstokens/accesstokensmock"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/globals"
	"github.com/calmisland/go-testify/mock"
)

// Setup Setup
func Setup() {
	if err := services.Initialize(services.ProductConfig{
		DownloadBaseURL: "http://localhost/",
	}); err != nil {
		panic(err)
	}

	setupProductDatabase()
	setupProductService()
	setupProductAccessService()
	setupPassService()
	setupPassAccessService()
	setupContentService()
	setupCORS()

	setupAccessTokenSystems()

	globals.Verify()
}

func setupProductDatabase() {
	globals.ProductDatabase = productmemorydb.New()
}

func setupProductService() {
	globals.ProductService = &productservice.StandardProductService{
		ProductDatabase: globals.ProductDatabase,
	}
}

func setupProductAccessService() {
	globals.ProductAccessService = &productaccessservice.StandardProductAccessService{
		ProductDatabase: globals.ProductDatabase,
	}
}

func setupPassService() {
	globals.PassService = &passservice.StandardPassService{
		ProductDatabase: globals.ProductDatabase,
	}
}

func setupPassAccessService() {
	globals.PassAccessService = &passaccessservice.StandardPassAccessService{
		ProductDatabase: globals.ProductDatabase,
	}
}

func setupContentService() {
	globals.ContentService = &contentservice.StandardContentService{
		ProductDatabase: globals.ProductDatabase,
	}
}

func setupCORS() {
	globals.CORSOptions = &apirouter.CORSOptions{
		AllowAnyOrigin: true,
		AllowedOrigins: []string{},
		AllowedMethods: []string{
			"DELETE",
			"GET",
			"OPTIONS",
			"POST",
			"PUT",
		},
		AllowCredentials: false,
		AllowedHeaders: []string{
			"Authorization",
			"Accept",
			"Accept-Encoding",
			"Content-Type",
		},
		ExposeAllHeaders:    true,
		ExposedHeaders:      []string{},
		AccessControlMaxAge: 7200,
	}
}

func setupAccessTokenSystems() {
	accessTokenValidator := &accesstokensmock.MockValidator{}
	accessTokenValidator.On("ValidateAccessToken", mock.Anything).Return(&sessions.SessionData{
		SessionID: "TEST-SESSION",
		AccountID: "TEST-ACCOUNT",
		DeviceID:  "TEST-DEVICE",
	}, nil)

	globals.AccessTokenValidator = accessTokenValidator
}
