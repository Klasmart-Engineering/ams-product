// +build !china

package server

import (
	"bitbucket.org/calmisland/go-server-aws/awsdynamodb"
	"bitbucket.org/calmisland/go-server-configs/configs"
	"bitbucket.org/calmisland/go-server-logs/errorreporter"
	"bitbucket.org/calmisland/go-server-logs/errorreporter/slackreporter"
	"bitbucket.org/calmisland/go-server-product/productdatabase/productdynamodb"
	"bitbucket.org/calmisland/go-server-requests/tokens/accesstokens"
	"bitbucket.org/calmisland/product-lambda-funcs/src/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/src/services"
)

// Setup setup the server based on configuration
func Setup() {
	if err := awsdynamodb.InitializeFromConfigs(); err != nil {
		panic(err)
	}
	if err := services.InitializeFromConfigs(); err != nil {
		panic(err)
	}

	setupAccessTokenSystems()
	setupLPAccessTokenSystems()

	productdynamodb.ActivateDatabase()

	setupSlackReporter()

	globals.Verify()
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

func setupLPAccessTokenSystems() {
	var validatorLPConfig accesstokens.ValidatorConfig
	err := configs.LoadConfig("learn_and_play_session", &validatorLPConfig, true)
	if err != nil {
		panic(err)
	}

	globals.AccessTokenLPValidator, err = accesstokens.NewLPValidator(validatorLPConfig)
	if err != nil {
		panic(err)
	}
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
