package router

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/hive/global"
	"github.com/beewit/hive/handler"

	"github.com/labstack/echo"

	"fmt"

	"github.com/beewit/beekit/utils/convert"
	"github.com/labstack/echo/middleware"
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
	e.POST("/api/rules/list", handler.GetRules, handler.Filter)
	e.POST("/api/func/list", handler.GetFuncList, handler.Filter)
	e.POST("/api/func/account/list", handler.GetFuncAllByIdsAndAccId, handler.Filter)
	e.POST("/api/func/account/funcId", handler.GetEffectiveFuncById, handler.Filter)
	e.POST("/api/func/account/funcList", handler.GetEffectiveFuncList, handler.Filter)
	e.POST("/api/func/account/group", handler.GetFuncGroupByAccountId, handler.Filter)
	e.POST("/api/func/account/share/wechat/app/add", handler.AddShareWechatAppTime, handler.Filter)

	e.POST("/api/account/auths", handler.GetAccountAuths, handler.Filter)
	e.POST("/api/account/func/list", handler.GetAccountFuncList, handler.Filter)
	e.POST("/api/account/updatePwd", handler.UpdatePassword, handler.Filter)
	e.POST("/api/account/wechat/group/add", handler.AddAccountWechatGroup, handler.Filter)
	e.POST("/api/account/action/log/add", handler.AddActionLogs, handler.Filter)
	e.POST("/api/account/identity/auth/add", handler.AddIdentityAuth, handler.Filter)
	e.POST("/api/account/identity/auth/get", handler.GetIndetityAuth, handler.Filter)
	e.POST("/api/account/func/give/log", handler.GetFuncGiveLog, handler.Filter)
	e.POST("/api/account/share/account/list", handler.GetShareAccountList, handler.Filter)

	e.POST("/api/order/pay/list", handler.GetPayOrderList, handler.Filter)
	e.POST("/api/wechat/group/list", handler.GetWechatGroupList, handler.Filter)
	e.POST("/api/wechat/group/class", handler.GetWechatGroupClass, handler.Filter)

	e.POST("/api/add/system/log", handler.AddSystemLog)
	utils.Open(global.Host)
	port := ":" + convert.ToString(global.Port)
	e.Logger.Fatal(e.Start(port))
}
