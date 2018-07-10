package controller

import (
	"github.com/labstack/echo"
	"net/http"
	"log"
	"ulang.com/wpvp/model"
)

func SysUserUpdateInfo(c echo.Context) error {
	curUser := GetCurrentSysUser(c)
	user := UserForm{}
	if err := c.Bind(user); err != nil {
		log.Println("========Bind error==========, err:", err)
		return c.JSON(http.StatusOK, ResponseMsg(S_FORM_ERROR))
	}

	if len(user.Password) < 6 {
		return c.JSON(http.StatusOK, ResponseMsg(S_SHORT_PASSWORD))
	}

	_, found := model.FindSysUser(user.Username)
	if !found {
		return c.JSON(http.StatusOK, ResponseMsg(S_INVALID_USERID))
	}
	curUser.Username = user.Username
	curUser.Password = user.Password
	curUser.PhoneNumber = user.PhoneNumber
	curUser.Avatar = user.Avatar
	curUser.Email = user.Email
	curUser.Update()

	return c.JSON(http.StatusOK, echo.Map{
		"Code":    S_OK,
		"Message": GetMessage(S_OK),
	})
}
func SysUserRegister(c echo.Context) error {
	username := c.FormValue("Username")
	password := c.FormValue("Password")
	if len(password) == 0 || len(username) == 0 {
		return c.JSON(http.StatusOK, ResponseMsg(S_INVALID_PARAM))
	}
	if len(password) < 6 {
		return c.JSON(http.StatusOK, ResponseMsg(S_SHORT_PASSWORD))
	}
	user, found := model.FindSysUser(username)
	if found {
		return c.JSON(http.StatusOK, ResponseMsg(S_DUP_USERID))
	}

	user, err := model.NewSysUser(username, password)
	if err != nil {
		return c.JSON(http.StatusOK, ResponseMsg(S_DATABASE_ERROR))
	}
	return c.JSON(http.StatusOK, echo.Map{
		"Code":    S_OK,
		"Message": GetMessage(S_OK),
		"Data": echo.Map{
			"UserId":   user.UserId,
			"Username": user.Username,
		},
	})
}
func SysUserLogin(c echo.Context) error {
	username := c.FormValue("Username")
	password := c.FormValue("Password")
	if len(password) == 0 || len(username) == 0 {
		return c.JSON(http.StatusOK, ResponseMsg(S_INVALID_PARAM))
	}
	if len(password) < 6 {
		return c.JSON(http.StatusOK, ResponseMsg(S_SHORT_PASSWORD))
	}
	user, found := model.FindSysUser(username)
	if !found {
		return c.JSON(http.StatusOK, ResponseMsg(S_DUP_USERID))
	}
	ok := user.Validate(password)
	if !ok {
		return c.JSON(http.StatusOK, ResponseMsg(S_INVALID_PASSWORD))
	}
	SetCurrentSysUser(c, user.Username)

	return c.JSON(http.StatusOK, echo.Map{
		"Code":    S_OK,
		"Message": GetMessage(S_OK),
		"Data": echo.Map{
			"UserId":   user.UserId,
			"Username": user.Username,
		},
	})
}
func SysUserLogout(c echo.Context) error {
	DeleteCurrentSysUser(c)

	return c.JSON(http.StatusOK, echo.Map{
		"Code":    S_OK,
		"Message": GetMessage(S_OK),
	})
}
