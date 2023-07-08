package core

import "github.com/labstack/echo/v4"

type Module interface {
	Init(g *echo.Group)
}

func Setup(e *echo.Echo, modules ...Module) {
	api := e.Group("/api")

	for _, module := range modules {
		module.Init(api)
	}
}
