package main

import (
	"WalletVerifyDemo/router"
	"WalletVerifyDemo/tools"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"os"
)

var (
	e *echo.Echo
)

func init() {
	e = echo.New()
	e.Use(middleware.Logger())
	e.Validator = tools.NewCustomerValidator()
	router.InitRoutes(e)
}

func getServePort() string {
	port := os.Getenv("WALLET_VERIFY_PORT")
	if port == "" {
		port = "1323"
	}

	return port
}

func main() {
	e.Logger.Fatal(e.Start("127.0.0.1:" + getServePort()))
}
