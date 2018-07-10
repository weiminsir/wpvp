package controller

import (
	"strings"
	"github.com/labstack/echo"
	"net/http"
	"ulang.com/wpvp/model"
	"github.com/labstack/echo-contrib/session"
)

const (
	SYSU_KEY = "sysuser_wpvp"
	API_KEY  = "api"
)

func SysUserAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if strings.Contains(c.Request().URL.Path, "admin/login") ||
				strings.Contains(c.Request().URL.Path, "admin/register") {
				return next(c)
			}
			sess, _ := session.Get("wpvp", c)
			if sess.Values[SYSU_KEY] == nil {
				return next(c)
			}
			username := sess.Values[SYSU_KEY].(string)
			if len(username) != 0 {
				user, ok := model.FindSysUser(username)
				if !ok {
					return c.JSON(http.StatusOK, ResponseMsg(S_LOGIN_REQUIRED))
				}
				c.Set(SYSU_KEY, user)
				return next(c)
			} else {
				return c.JSON(http.StatusOK, ResponseMsg(S_LOGIN_REQUIRED))
			}
		}
	}
}
func SetCurrentSysUser(c echo.Context, username string) {
	sess, _ := session.Get("wpvp", c)
	sess.Values[SYSU_KEY] = username
	sess.Save(c.Request(), c.Response())
}

func GetCurrentSysUser(c echo.Context) (*model.SysUser) {
	return c.Get(SYSU_KEY).(*model.SysUser)
}
func DeleteCurrentSysUser(c echo.Context) {
	sess, _ := session.Get("wpvp", c)
	delete(sess.Values, SYSU_KEY)
	sess.Save(c.Request(), c.Response())
}
