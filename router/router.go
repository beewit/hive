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
	e.Use(middleware.Recover())

	e.Static("/app", "app")
	e.File("/", "app/page/index.html")

	e.POST("/api/template", handler.GetTemplateByListPage, handler.Filter)
	e.POST("/api/template/update/refer/:id", handler.UpdateTemplateReferById, handler.Filter)
	e.POST("/api/template/:id", handler.GetTemplateById, handler.Filter)
	e.POST("/api/platform", handler.GetPlatformList, handler.Filter)
	e.POST("/api/platform/one", handler.GetPlatformId, handler.Filter)
	e.POST("/api/rules/list",handler.GetRules,handler.Filter)

	utils.Open(global.Host)

	port := ":" + convert.ToString(global.Port)

	e.Logger.Fatal(e.Start(port))
}
