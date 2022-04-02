package user

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/wrpota/go-echo/internal/pkg/response"
)

type User struct {
}

type UserRequest struct {
	Name  string `json:"name" form:"name" query:"name"`
	Email string `json:"email" form:"email" query:"email"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) HelloWorld(c echo.Context) error {
	fmt.Println("Hello World")
	user := new(UserRequest)
	if err := c.Bind(user); err != nil {
		return nil
	}
	return response.Success(c, "success", user)
}
