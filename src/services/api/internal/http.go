package internal

import (
	"context"
	"github.com/artbred/aliasflux/src/pkg/common"
	"github.com/artbred/aliasflux/src/pkg/config"
	"github.com/artbred/aliasflux/src/services/api/internal/monitoring"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Setup(e *echo.Echo) {
	SetupDocs(e)
	SetupValidator(e)
	SetupMiddleware(e)
	monitoring.StartMonitoring(e)
}

func ServerHttp(e *echo.Echo) {
	addr := ":" + config.ApiPort

	if config.Debug {
		err := e.Start(addr)
		if err != nil {
			common.Logger.Fatal(err)
		}
	} else {
		go func() {
			if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
				common.Logger.Fatal("shutting down the server")
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)
		<-quit

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			common.Logger.Fatal(err)
		}
	}
}
