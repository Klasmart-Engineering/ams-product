// +build !china

package server

import (
	"bitbucket.org/calmisland/go-server-shared/security"
	"bitbucket.org/calmisland/go-server-shared/services/aws/awsdynamodb"
	"bitbucket.org/calmisland/go-server-standard/databases/productdatabase/productdynamodb"
)

// Setup setup the server based on configuration
func Setup() {
	if err := security.InitializeFromConfigs(); err != nil {
		panic(err)
	}
	if err := awsdynamodb.InitializeFromConfigs(); err != nil {
		panic(err)
	}

	productdynamodb.ActivateDatabase()
}
