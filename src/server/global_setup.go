// +build !china

package server

import (
	"bitbucket.org/calmisland/go-server-shared/v2/configs"
	"bitbucket.org/calmisland/go-server-shared/v2/errors/errorreporter"
	"bitbucket.org/calmisland/go-server-shared/v2/errors/errorreporter/slackreporter"
	"bitbucket.org/calmisland/go-server-shared/v2/security"
	"bitbucket.org/calmisland/go-server-shared/v2/services/aws/awsdynamodb"
	"bitbucket.org/calmisland/go-server-standard/databases/productdatabase/productdynamodb"
	"bitbucket.org/calmisland/product-lambda-funcs/src/services"
)

// Setup setup the server based on configuration
func Setup() {
	if err := security.InitializeFromConfigs(); err != nil {
		panic(err)
	}
	if err := awsdynamodb.InitializeFromConfigs(); err != nil {
		panic(err)
	}
	if err := services.InitializeFromConfigs(); err != nil {
		panic(err)
	}

	productdynamodb.ActivateDatabase()

	setupSlackReporter()
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
