package router

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"ulang.com/wpvp/controller"
	"ulang.com/wpvp/conf"
	"net/url"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
)

var e *echo.Echo
/**
* echo 初始化设置
 */
func Start() {
	e = echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	store := sessions.NewCookieStore([]byte("ulang123456!@#$%^"))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	e.Use(session.Middleware(store))
	e.Static("/", "./static/index.html")
	e.Static("/static", "./static")
	initRouter()
	e.Logger.Fatal(e.Start(conf.Config.Port))
}

func initRouter() {
	adminRouter := e.Group("admin")
	adminRouter.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE,echo.OPTIONS},
	}))
	adminRouter.Use(controller.SysUserAuth())
	{
		adminRouter.GET("/speakers/:task_id", controller.AdminGetSpkRecords)
		adminRouter.GET("/tasks", controller.AdminGetRetrievalTasks)
		adminRouter.GET("/tasks/:page_start/:page_size", controller.AdminGetRetrievalTasksList)
		adminRouter.POST("/tasks", controller.AdminInsertRetrievalTask)
		adminRouter.POST("/register", controller.SysUserRegister)
		adminRouter.PUT("/user", controller.SysUserUpdateInfo)
		adminRouter.POST("/login", controller.SysUserLogin)
		adminRouter.POST("/logout", controller.SysUserLogout)
	}
	APIRouter := e.Group("/api")
	APIRouter.Use(controller.SysUserAuth())
	{
		url, err := url.Parse(conf.Config.Proxy_Addr)
		if err != nil {
			e.Logger.Fatal(err)
		}
		APIRouter.Use(middleware.Proxy(&middleware.RoundRobinBalancer{
			Targets: []*middleware.ProxyTarget{
				{
					URL: url,
				},
			},
		}))
	}

}
