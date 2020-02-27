package globalsetup

import (
	"bitbucket.org/calmisland/go-server-aws/awsdynamodb"
	"bitbucket.org/calmisland/go-server-configs/configs"
	"bitbucket.org/calmisland/go-server-logs/errorreporter"
	"bitbucket.org/calmisland/go-server-logs/errorreporter/slackreporter"
	"bitbucket.org/calmisland/go-server-product/contentservice"
	"bitbucket.org/calmisland/go-server-product/eventticketregionservice"
	"bitbucket.org/calmisland/go-server-product/passaccessservice"
	"bitbucket.org/calmisland/go-server-product/passservice"
	"bitbucket.org/calmisland/go-server-product/productaccessservice"
	"bitbucket.org/calmisland/go-server-product/productdatabase/productdynamodb"
	"bitbucket.org/calmisland/go-server-product/productservice"
	"bitbucket.org/calmisland/go-server-requests/apirouter"
	"bitbucket.org/calmisland/go-server-requests/tokens/accesstokens"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/src/services"
)

// Setup Setup
func Setup() {
	// Setup the Slack reporter first
	setupSlackReporter()

	if err := services.InitializeFromConfigs(); err != nil {
		panic(err)
	}

	setupProductDatabase()
	setupProductService()
	setupProductAccessService()
	setupPassService()
	setupPassAccessService()
	setupContentService()
	setupEventTicketRegionService()

	setupAccessTokenSystems()
	setupCORS()

	globals.Verify()
}

func setupProductDatabase() {
	var productDatabaseConfig awsdynamodb.ClientConfig
	err := configs.LoadConfig("product_database_dynamodb", &productDatabaseConfig, true)
	if err != nil {
		panic(err)
	}

	ddbClient, err := awsdynamodb.NewClient(&productDatabaseConfig)
	if err != nil {
		panic(err)
	}

	globals.ProductDatabase, err = productdynamodb.New(ddbClient)
	if err != nil {
		panic(err)
	}
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

func setupEventTicketRegionService() {
	globals.EventTicketRegionService = &eventticketregionservice.StandardService{
		ProductDatabase: globals.ProductDatabase,
	}
}

func setupAccessTokenSystems() {
	var validatorConfig accesstokens.ValidatorConfig
	err := configs.LoadConfig("access_tokens", &validatorConfig, true)
	if err != nil {
		panic(err)
	}
	globals.AccessTokenValidator, err = accesstokens.NewValidator(validatorConfig)
	if err != nil {
		panic(err)
	}
}

func setupCORS() {
	var corsConfig apirouter.CORSOptions
	err := configs.LoadConfig("cross_origin_resource_sharing", &corsConfig, true)
	if err != nil {
		panic(err)
	}

	globals.CORSOptions = &corsConfig
}

func setupSlackReporter() {
	var slackReporterConfig slackreporter.Config
	err := configs.LoadConfig("error_reporter_slack", &slackReporterConfig, false)
	if err != nil {
		panic(err)
	}

	// Check if there is a configuration for the Slack error reporter
	if len(slackReporterConfig.HookURL) == 0 {
		return
	}

	reporter, err := slackreporter.New(&slackReporterConfig)
	if err != nil {
		panic(err)
	}

	errorreporter.Active = reporter
}
