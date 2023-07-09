package chats

import (
	"github.com/artbred/aliasflux/src/pkg/config"
	"github.com/artbred/aliasflux/src/services/api/internal"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Router struct {
	BasePath string
	upgrader websocket.Upgrader
}

func (r *Router) Init(g *echo.Group) {
	g = g.Group(r.BasePath)

	g.Use(internal.RateLimitMiddleware())

	g.GET("/:id", r.GetChatHandler)
	g.POST("/create", r.CreateChatHandler)
	g.GET("/ws", r.WebsocketHandler)
	g.GET("/settings", r.ListAvailableChatSettingsHandler)
}

func NewRouter() *Router {
	return &Router{
		BasePath: "/chats",
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return config.Debug
			},
		},
	}
}
