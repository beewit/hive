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
	e.File("/.well-known/pki-validation/fileauth.txt", "app/fileauth.txt")
	e.File("8VHBgcXdwx.txt", "8VHBgcXdwx.txt")
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
	e.POST("/api/account/identity/auth/get", handler.GetIdentityAuth, handler.Filter)
	e.POST("/api/account/func/give/log", handler.GetFuncGiveLog, handler.Filter)
	e.POST("/api/account/share/account/list", handler.GetShareAccountList, handler.Filter)
	e.POST("/api/account/advert/list", handler.GetAccountAdvertList, handler.Filter)
	e.POST("/api/account/advert/add", handler.AddAccountAdvert, handler.Filter)
	e.POST("/api/account/advert/get", handler.GetAccountAdvert, handler.Filter)
	e.POST("/api/account/advert/del", handler.DelAccountAdvert, handler.Filter)

	e.POST("/api/account/commpany/auth/submit", handler.SubmitAccountCompanyAuth, handler.Filter)
	e.POST("/api/account/commpany/auth/get", handler.GetAccountCompanyAuthByID, handler.Filter)
	e.POST("/api/account/redpacket/card/add", handler.AddRedPacketCard, handler.Filter)
	e.POST("/api/account/redpacket/card/del", handler.DeleteRedPacketCard, handler.Filter)
	e.POST("/api/account/redpacket/card/get", handler.GetRedPacketCardByID, handler.Filter)
	e.POST("/api/account/redpacket/card/list", handler.GetRedPacketCardList, handler.Filter)
	e.POST("/api/account/redpacket/card/def", handler.GetRedPacketCardDef, handler.Filter)

	e.POST("/api/app/setting", handler.GetAppSetting)

	e.POST("/api/account/wallet", handler.GetWallet, handler.Filter)
	e.POST("/api/account/withdrawCash/card/add", handler.AddWithdrawCashCard, handler.Filter)
	e.POST("/api/account/withdrawCash/apply", handler.ApplyWithdrawCash, handler.Filter)
	e.POST("/api/account/withdrawCash/apply/list", handler.GetApplyWithdrawCashList, handler.Filter)
	e.POST("/api/account/wallet/log", handler.GetWalletLogList, handler.Filter)

	e.POST("/api/order/pay/list", handler.GetPayOrderList, handler.Filter)
	e.POST("/api/wechat/group/list", handler.GetWechatGroupList, handler.Filter)
	e.POST("/api/wechat/group/class", handler.GetWechatGroupClass, handler.Filter)

	e.POST("/api/add/system/log", handler.AddSystemLog)

	// --- 功能试用期 ---
	e.POST("/api/account/func/tryout/get", handler.GetAccountFuncTryout, handler.Filter)
	e.POST("/api/account/func/tryout/add", handler.AddAccountFuncTryout, handler.Filter)

	// --- 移动端自动营销发送记录 ---
	e.POST("/api/account/func/handle/log/add", handler.AddAccountFuncHandleLog, handler.Filter)
	e.POST("/api/account/func/handle/log/group", handler.GetAccountFuncHandleGroup, handler.Filter)
	e.POST("/api/account/func/handle/log/list", handler.GetAccountFuncHandleLogList, handler.Filter)

	//小程序
	//现金券
	e.POST("/api/account/coupon/add", handler.AddCoupon, handler.Filter)
	e.POST("/api/account/coupon/use", handler.UseCoupon, handler.Filter)
	e.POST("/api/account/coupon/delete", handler.DeleteCoupon, handler.Filter)
	//发送现金券列表
	e.POST("/api/account/coupon/list", handler.GetCouponList, handler.Filter)
	e.POST("/api/account/coupon/use/num", handler.GetUseCouponNum, handler.Filter)
	//现金券详情
	e.POST("/api/account/coupon/get", handler.GetCouponByID, handler.Filter)
	//红包
	e.POST("/api/account/redpacket", handler.GetRedPacketById)
	//发红包
	e.POST("/api/account/send/redpacket", handler.SendRedPacket, handler.Filter)
	e.POST("/api/account/send/redpacket/sumprice", handler.GetSendRedPacketSumPrice, handler.Filter)
	//发红包记录
	e.POST("/api/account/send/redpacket/list", handler.GetSendRedPacketList, handler.Filter)

	e.POST("/api/account/redpacket/access/log/add", handler.AddRedPacketAccessLog)
	e.POST("/api/account/redpacket/access/log/num", handler.GetRedPacketAccessLogNum, handler.Filter)

	//[小程序非系统账号也可访问，根据小程序用户标识]
	//领取现金券
	e.POST("/api/account/receive/coupon", handler.ReceiveCoupon)
	e.POST("/api/account/receive/coupon/qrcode", handler.CreateCouponQrCode)
	//红包的领取记录和领取现金券记录
	e.POST("/api/account/redpacket/receive/coupon/record", handler.GetReceiveRedPacketAndCouponList)
	//领取的红包记录
	e.POST("/api/account/redpacket/receive/list", handler.GetReceiveRedPacketList)
	//创建红包领取记录
	e.POST("/api/account/redpacket/receive", handler.ReceiveRedPacket)
	//领取现金券记录
	e.POST("/api/account/receive/coupon/list", handler.GetReceiveCouponList)
	//分享红包
	e.POST("/api/account/share/redpacket", handler.AddShareRedPacket)
	//分享红包
	e.POST("/api/account/share/redpacket/num", handler.GetShareRedPacketCountByRedPacketId)

	//帮助
	e.POST("/api/help/list", handler.GetHelpList)
	e.GET("/api/help/list", handler.GetHelpList)
	e.POST("/api/help/get", handler.GetHelp)
	e.GET("/api/help/get", handler.GetHelp)

	//地区
	e.POST("/api/area/get/child", handler.GetAreaChild)
	e.POST("/api/area/get", handler.GetArea)
	e.GET("/api/area/get/child", handler.GetAreaChild)
	e.GET("/api/area/get", handler.GetArea)

	utils.Open(global.Host)
	port := ":" + convert.ToString(global.Port)
	e.Logger.Fatal(e.Start(port))
}
