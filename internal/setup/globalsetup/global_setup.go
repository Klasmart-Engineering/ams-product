package globalsetup

import (
	"errors"
	"fmt"
	"os"

	"bitbucket.org/calmisland/go-server-aws/awsdynamodb"
	"bitbucket.org/calmisland/go-server-configs/configs"
	"bitbucket.org/calmisland/go-server-logs/errorreporter"
	"bitbucket.org/calmisland/go-server-logs/errorreporter/slackreporter"
	"bitbucket.org/calmisland/go-server-product/contentservice"
	"bitbucket.org/calmisland/go-server-product/klppassservice"
	"bitbucket.org/calmisland/go-server-product/passaccessservice"
	"bitbucket.org/calmisland/go-server-product/passservice"
	"bitbucket.org/calmisland/go-server-product/productaccessservice"
	"bitbucket.org/calmisland/go-server-product/productdatabase/productdynamodb"
	"bitbucket.org/calmisland/go-server-product/productservice"
	"bitbucket.org/calmisland/go-server-requests/tokens/accesstokens"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/globals"
	v1Services "bitbucket.org/calmisland/product-lambda-funcs/internal/services/v1"
	"github.com/getsentry/sentry-go"
)

// Setup Setup
func Setup() {
	setupSentry()
	// Setup the Slack reporter first
	setupSlackReporter()

	if err := v1Services.InitializeFromConfigs(); err != nil {
		panic(err)
	}

	setupProductDatabase()
	setupProductService()
	setupProductAccessService()
	setupPassService()
	setupPassAccessService()
	setupContentService()
	setupKlpPassService()

	setupAccessTokenSystems()

	globals.Verify()
}

func setupSentry() {
	var env string = fmt.Sprintf("%s", os.Getenv("SERVER_STAGE"))

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:         "https://619f73d7336e479ba6c1659cd256b4f8@o412774.ingest.sentry.io/5614058",
		Environment: env,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}
}

func setupProductDatabase() {
	var productDatabaseConfig awsdynamodb.ClientConfig
	err := configs.ReadEnvConfig(&productDatabaseConfig)
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

func setupKlpPassService() {
	globals.KlpPassService = &klppassservice.StandardKlpPassService{
		ProductDatabase: globals.ProductDatabase,
	}
}

func setupAccessTokenSystems() {
	var err error
	var validatorConfig accesstokens.ValidatorConfig

	bPublicKey := configs.LoadBinary("account.pub")
	if bPublicKey == nil {
		panic(errors.New("the account.pub file is mandatory"))
	}

	validatorConfig.PublicKey = string(bPublicKey)

	globals.AccessTokenValidator, err = accesstokens.NewValidator(validatorConfig)
	if err != nil {
		panic(err)
	}
}

func setupSlackReporter() {
	var slackReporterConfig slackreporter.Config
	err := configs.ReadEnvConfig(&slackReporterConfig)
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
