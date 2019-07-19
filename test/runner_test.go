package test

import (
	"testing"

	"bitbucket.org/calmisland/go-server-configs/configs"
	"bitbucket.org/calmisland/go-server-product/testproductdatabase"
	"bitbucket.org/calmisland/go-server-product/testproductdatabase/testproductdynamodb"
	"bitbucket.org/calmisland/product-lambda-funcs/src/server"
)

// TestProduct execute a defined list of tests, using the project configuration environment
func TestProduct(t *testing.T) {
	err := configs.UpdateConfigDirectoryPath("../" + configs.DefaultConfigFolderName)
	if err != nil {
		panic(err)
	}
	server.Setup()
	testproductdatabase.RunTestSuite(t, new(testproductdynamodb.ProductDynamoDBClient))
}
