package user

import (
	"github.com/labstack/echo/v4"
	"github.com/wrpota/go-echo/internal/pkg/response"
)

type UserRequest struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

func (u *handler) HelloWorld() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := new(UserRequest)
		if err := c.Bind(user); err != nil {
			return nil
		}
		return response.Success(c, "success", user)
	}
}
