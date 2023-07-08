package internal

import (
	"fmt"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/pkg/config"
	"github.com/artbred/aliasflux/src/services/api/internal/monitoring"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"strconv"
)

func defaultSkipper(c echo.Context) bool {
	if c.Path() == "/metrics" || c.Path() == "/health" || c.Path() == "/swagger/*" {
		return true
	}

	return false
}

func RateLimitMiddleware() echo.MiddlewareFunc {
	rateLimiterConfig := middleware.DefaultRateLimiterConfig
	rateLimiterConfig.Store = middleware.NewRateLimiterMemoryStore(10)
	return middleware.RateLimiterWithConfig(rateLimiterConfig)
}

func ErrorMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)

			status := c.Response().Status
			path := c.Request().URL.Path
			method := c.Request().Method

			go func() {
				if status < 200 || status > 299 {
					monitoring.EndpointErrors.WithLabelValues(path, method, strconv.Itoa(status)).Inc()
				}
			}()

			return err
		}
	}
}

func SetupMiddleware(e *echo.Echo) {
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output:  common.Logger.Writer(),
		Skipper: defaultSkipper,
	}))

	corsConfig := middleware.DefaultCORSConfig
	if !config.Debug {
		corsConfig = middleware.CORSConfig{
			Skipper:          middleware.DefaultSkipper,
			AllowOrigins:     []string{fmt.Sprintf("https://%s", config.Domain), fmt.Sprintf("https://www.%s", config.Domain)},
			AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
			AllowHeaders:     []string{"*"},
			AllowCredentials: true,
		}
	}

	e.Use(middleware.CORSWithConfig(corsConfig))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.RequestID())

	e.Use(ErrorMiddleware())
}
