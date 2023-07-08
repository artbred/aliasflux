package main

import (
	"github.com/artbred/aliasflux/src/services/api/core"
	v1 "github.com/artbred/aliasflux/src/services/api/core/v1"
	"github.com/artbred/aliasflux/src/services/api/cron"
	"github.com/artbred/aliasflux/src/services/api/internal"
	"github.com/labstack/echo/v4"
	_ "go.uber.org/automaxprocs"
	"net/http"
)

func SetupAPI(e *echo.Echo) {
	core.Setup(e, v1.NewVersion())
}

// @title AliasFlux API
// @version 1.0

// @BasePath /api/v1
func main() {
	e := echo.New()
	internal.Setup(e)

	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	go cron.Start()

	SetupAPI(e)
	internal.ServerHttp(e)
}
