package user

import (
	"github.com/labstack/echo/v4"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()
	HelloWorld() echo.HandlerFunc
}

type handler struct {
}

func (h *handler) i() {}

//返回handler 对象
func New() Handler {
	return &handler{}
}
