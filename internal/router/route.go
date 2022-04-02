package router

import (
	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	"github.com/wrpota/go-echo/internal/api/user"
	"github.com/wrpota/go-echo/internal/global/variable"
	"github.com/wrpota/go-echo/internal/pkg/echozap"
)

func InitRouter() *echo.Echo {
	router := echo.New()
	// //根据配置进行设置跨域
	if variable.Config.GetBool("HttpServer.AllowCrossDomain") {
		router.Use(middleware.CORS())
	}
	router.Use(middleware.Recover())
	router.Use(echozap.ZapLogger(variable.EchoZapLog))

	var userRouter = router.Group("user")
	user := user.NewUser()
	userRouter.GET("", user.HelloWorld)

	return router
}
