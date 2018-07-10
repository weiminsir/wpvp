package filter

import (
	"github.com/labstack/echo"
)

/**
此方法为公共拦截器
这里做是否登录、权限校验等操作
给全局变量设值可以在Controller 中获取到
 */

func OpenTracing() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}
