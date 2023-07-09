package users

import (
	"github.com/artbred/aliasflux/src/services/api/internal"
	"github.com/labstack/echo/v4"
)

type Router struct {
	BasePath string
}

func (r *Router) Init(g *echo.Group) {
	g = g.Group(r.BasePath)

	g.Use(internal.RateLimitMiddleware())

	g.GET("/create", r.CreateUserHandler)
}

func NewRouter() *Router {
	return &Router{
		BasePath: "/users",
	}
}
