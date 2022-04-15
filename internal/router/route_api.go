package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wrpota/go-echo/internal/api/user"
)

func setApiRouter(r *echo.Echo) {
	// r.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	variable.EchoZapLog.Info("reqBody:" + string(reqBody))
	// }))

	api := r.Group("/api")
	var userRouter = api.Group("/user")
	user := user.New()
	userRouter.GET("", user.HelloWorld())
}
