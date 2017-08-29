package router

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/hive/handler"

	"github.com/labstack/echo"

	"fmt"
	"github.com/labstack/echo/middleware"
)

func Start() {
	fmt.Printf("登陆授权系统启动")

	e := echo.New()
	e.Use(middleware.Gzip())
	//e.Use(middleware.Logger())
	//e.Use(middleware.CSRF())
	//e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Static("/static", "static")
	e.Static("/page", "page")
	e.File("/", "page/login.html")

	e.GET("/api/template", handler.GetTemplateByList)
	e.POST("/api/template/:id", handler.GetTemplateById)

	utils.Open(global.Host)

	e.Logger.Fatal(e.Start(":8090"))
}
