package v1

import (
	"github.com/artbred/aliasflux/src/services/api/core/v1/chats"
	"github.com/artbred/aliasflux/src/services/api/core/v1/users"
	"github.com/labstack/echo/v4"
)

type Version struct {
	BasePath string
}

func (r *Version) Init(g *echo.Group) {
	g = g.Group(r.BasePath)

	chats.NewRouter().Init(g)
	users.NewRouter().Init(g)
}

func NewVersion() *Version {
	return &Version{
		BasePath: "/v1",
	}
}
