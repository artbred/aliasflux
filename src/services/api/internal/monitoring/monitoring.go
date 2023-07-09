package monitoring

import (
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	EndpointErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "api_endpoint_errors_total",
		Help: "Total number of failure support from api",
	}, []string{"endpoint", "method", "status_code"})
)

func StartMonitoring(e *echo.Echo) {
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
}
