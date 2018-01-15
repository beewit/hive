/*
领取红包记录
author：zxb
2018-01-13
*/

package handler

import (
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
)

func ReceiveRedPacket(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	redPacket := GetRedPacket(convert.MustInt64(id))
	if redPacket != nil {
		return utils.ErrorNull(c, "该红包不存在")
	}
	sendRedPacketAcc := GetAccountById(convert.MustInt64(redPacket["account_id"]))
	if sendRedPacketAcc == nil {
		return utils.ErrorNull(c, "发红包的账号已被锁定无法领取")
	}
	if redPacket["status"] != enum.NORMAL {
		return utils.ErrorNull(c, "该红包已失效")
	}
	if convert.MustFloat64(redPacket["money"])-convert.MustFloat64(redPacket["send_money"]) < convert.MustFloat64(redPacket["random_money"]) {
		return utils.ErrorNull(c, "该红包已领完")
	}
	println(acc.ID)
	// _, err = global.DB.InsertMap("account_receive_red_packet", map[string]interface{}{
	// 	"id":                         utils.ID,
	// 	"account_id":                 acc.ID,
	// 	"account_send_red_packet_id": redPacket["id"],
	// 	"ct_time":                    utils.CurrentTime(),
	// 	"ip":                         c.RealIP(),
	// 	"status":                     "未使用",
	// })
	// if err != nil {
	// 	return utils.ErrorNull(c, "领取红包失败")
	// }
	return utils.SuccessNull(c, "领取红包成功")
}

/**
*领取红包记录
 */
func GetReceiveRedPacketList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "v_account_receive_red_packet",
		Where:     "rp_account_id=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
	}, acc.ID)
	if err != nil {
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}
