package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func ReturnJson(Context echo.Context, httpCode int, dataCode int, msg string, data interface{}) error {
	return Context.JSON(httpCode, map[string]interface{}{
		"code": dataCode,
		"msg":  msg,
		"data": data,
	})
}

//ReturnJsonFromString 将json字符窜以标准json格式返回（例如，从redis读取json格式的字符串，返回给浏览器json格式）
func ReturnJsonFromString(Context echo.Context, httpCode int, jsonStr string) error {
	Context.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return Context.String(httpCode, jsonStr)
}

// 语法糖函数封装

//Success 直接返回成功
func Success(c echo.Context, msg string, data interface{}) error {
	return ReturnJson(c, http.StatusOK, http.StatusOK, msg, data)
}

//Fail 失败的业务逻辑
func Fail(c echo.Context, dataCode int, msg string, data interface{}) error {
	return ReturnJson(c, http.StatusBadRequest, dataCode, msg, data)
}
