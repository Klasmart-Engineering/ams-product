package test

import (
	"testing"

	"bitbucket.org/calmisland/go-server-shared/v2/configs"
	"bitbucket.org/calmisland/go-server-standard/test/testdatabases/testdynamodb"
	"bitbucket.org/calmisland/product-lambda-funcs/src/server"
	"github.com/calmisland/go-testify/suite"
)

// TestProduct execute a defined list of tests, using the project configuration environment
func TestProduct(t *testing.T) {
	err := configs.UpdateConfigDirectoryPath("../" + configs.DefaultConfigFolderName)
	if err != nil {
		panic(err)
	}
	server.Setup()
	suite.Run(t, new(testdynamodb.ProductTestSuite))
}
