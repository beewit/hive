package router

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/hive/handler"

	"github.com/labstack/echo"

	"fmt"

	"github.com/labstack/echo/middleware"
	"github.com/beewit/beekit/utils/convert"
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

	e.POST("/api/template", handler.GetTemplateByListPage)
	e.POST("/api/template/update/refer/:id", handler.UpdateTemplateReferById)
	e.POST("/api/template/:id", handler.GetTemplateById)
	e.POST("/api/platform", handler.GetPlatformList)
	e.POST("/api/platform/one", handler.GetPlatformId)

	utils.Open(global.Host)

	port := ":" + convert.ToString(global.Port)

	e.Logger.Fatal(e.Start(port))
}
