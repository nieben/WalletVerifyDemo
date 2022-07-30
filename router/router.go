package router

import (
	"WalletVerifyDemo/handler"
	"github.com/labstack/echo/v4"
)

var verifyHandler = new(handler.VerifyHandler)

func InitRoutes(e *echo.Echo) {
	gets := map[string]echo.HandlerFunc{
		"/get_message": verifyHandler.Message,
	}

	posts := map[string]echo.HandlerFunc{
		"/verify": verifyHandler.Verify,
	}

	for key, value := range gets {
		e.GET(key, value)
	}
	for key, value := range posts {
		e.POST(key, value)
	}
}
