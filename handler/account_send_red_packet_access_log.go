/*
红包访问日志
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
	"strings"
)

func AddRedPacketAccessLog(c echo.Context) error {
	ws := GetOauthUser(c)
	if ws == nil {
		return utils.AuthWechatFailNull(c)
	}
	id := strings.TrimSpace(c.FormValue("id"))
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	redPacket := GetRedPacket(convert.MustInt64(id))
	if redPacket == nil {
		return utils.ErrorNull(c, "红包不存在或已过期")
	}
	acc := GetAccountByUnionId(ws.UnionId, enum.WECHAT)
	var accId interface{}
	if acc != nil {
		accId = acc["id"]
	} else {
		accId = nil
	}
	_, err := global.DB.InsertMap("account_send_red_packet_access_log", map[string]interface{}{
		"id": utils.ID(),
		"account_send_red_packet_id": id,
		"account_id":                 accId,
		"ct_time":                    utils.CurrentTime(),
		"wx_union_id":                ws.UnionId,
		"ip":                         c.RealIP(),
	})
	if err != nil {
		global.Log.Error("AddShareRedPacket account_send_red_packet_access_log sql error:%s", err.Error())
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	return utils.SuccessNull(c, "添加红包访问记录成功")
}

func GetRedPacketAccessLogNum(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	rows, err := global.DB.Query("SELECT count(1) as num FROM v_account_send_red_packet_access_log WHERE redPacketAccountId=? ", acc.ID)
	if err != nil {
		global.Log.Error("v_account_send_red_packet_access_log sql error:%s", err.Error())
		return utils.ErrorNull(c, "数据异常，"+err.Error())
	}
	if len(rows) != 1 {
		return utils.SuccessNullMsg(c, 0)
	}
	num := convert.MustInt64(rows[0]["num"])
	if convert.MustInt64(rows[0]["num"]) <= 0 {
		return utils.SuccessNullMsg(c, 0)
	}
	return utils.SuccessNullMsg(c, num)
}
