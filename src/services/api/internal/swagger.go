package internal

import (
	"github.com/artbred/aliasflux/src/pkg/config"
	_ "github.com/artbred/aliasflux/src/services/api/docs"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SetupDocs(e *echo.Echo) {
	if config.Debug {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}
}
