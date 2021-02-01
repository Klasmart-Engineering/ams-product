// +build !lambda

package main

import (
	"bitbucket.org/calmisland/go-server-configs/configs"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/globals"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/routers"
	"bitbucket.org/calmisland/product-lambda-funcs/internal/setup/globalsetup"
	"github.com/labstack/echo"
)

func main() {
	err := configs.UpdateConfigDirectoryPath(configs.DefaultConfigFolderName)
	if err != nil {
		panic(err)
	}

	globalsetup.Setup()

	echo := routers.SetupRouter()

	// Start server
	echo.Logger.Fatal(echo.Start(":8044"))

}

func createTablesRequest(c echo.Context) error {
	err := globals.ProductDatabase.CreateDatabaseTables()
	if err != nil {
		return err
	}
	return nil
}
