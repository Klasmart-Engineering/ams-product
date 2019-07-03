// +build !china

package server

import (
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
}
