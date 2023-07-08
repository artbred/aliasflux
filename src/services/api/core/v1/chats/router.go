package chats

import (
	"github.com/artbred/aliasflux/src/services/api/internal"
	"github.com/labstack/echo/v4"
)

type Router struct {
	BasePath string
}

func (r *Router) Init(g *echo.Group) {
	g = g.Group(r.BasePath)

	g.POST("/create", r.CreateChatHandler, internal.RateLimitMiddleware())
}

func NewRouter() *Router {
	return &Router{
		BasePath: "/chats",
	}
}
