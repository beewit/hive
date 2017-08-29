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
	e.Use(middleware.CSRF())
	e.Use(middleware.CORS())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())

	e.Static("/static", "static")
	e.Static("/page", "page")
	e.File("/", "page/login.html")

	e.GET("/api/template", handler.GetTemplateByList)
	e.GET("/api/template/id", handler.GetTemplateById)

	e.GET("test", jsonTest)

	utils.Open(global.Host)

	e.Logger.Fatal(e.Start(":8090"))
}

func jsonTest(c echo.Context) error {

	maps2 := make(map[string]interface{})
	maps2["id"] = 1
	maps2["name"] = "张三"
	maps := make(map[string]interface{})
	maps["id"] = 1
	maps["name"] = "张三"

	maps2["list"]=maps
	return  utils.Success(c, "", maps2)
}
