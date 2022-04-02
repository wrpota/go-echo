package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/wrpota/go-echo/configs"
	"github.com/wrpota/go-echo/internal/global/variable"
	"github.com/wrpota/go-echo/internal/pkg/echozap"
	"golang.org/x/time/rate"
)

func NewHttpServer() *echo.Echo {
	router := echo.New()
	//根据配置进行设置跨域
	if configs.Get().GetBool("HttpServer.AllowCrossDomain") {
		router.Use(middleware.CORS())
	}
	//流控
	if limit := configs.Get().GetFloat64("HttpServer.MaxRequestsPerSecond"); limit > 0 {
		configs := middleware.RateLimiterConfig{
			Store: middleware.NewRateLimiterMemoryStoreWithConfig(
				middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(limit)},
			),
			DenyHandler: func(context echo.Context, identifier string, err error) error {
				return context.JSON(http.StatusTooManyRequests, nil)
			},
		}

		router.Use(middleware.RateLimiterWithConfig(configs))
	}

	router.Use(middleware.Recover())
	router.Use(echozap.ZapLogger(variable.EchoZapLog))

	setApiRouter(router)

	return router
}
