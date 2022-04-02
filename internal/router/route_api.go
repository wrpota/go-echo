package router

import (
	"github.com/labstack/echo/v4"
	"github.com/wrpota/go-echo/internal/api/user"
)

func setApiRouter(r *echo.Echo) {
	api := r.Group("/api")
	var userRouter = api.Group("/user")
	user := user.NewUser()
	userRouter.GET("", user.HelloWorld)
}
