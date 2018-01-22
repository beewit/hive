/*
发送红包记录
author：zxb
2018-01-13
*/

package handler

import (
	"fmt"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
	"strings"
)

func GetRedPacketById(c echo.Context) error {
	ws, err := GetMiniAppSession(c)
	if err != nil || ws == nil {
		return utils.AuthFailNull(c)
	}
	id := strings.TrimSpace(c.FormValue("id"))
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	redPacketId := convert.MustInt64(id)
	redPacket := GetRedPacket(redPacketId)
	if redPacket == nil {
		return utils.ErrorNull(c, "红包不存在或已过期")
	}
	redPacket["account_id"] = nil
	var couponMaps []map[string]interface{}
	joinCouponIds := convert.ToString(redPacket["join_coupon_ids"])
	if joinCouponIds != "" {
		couponMaps = GetCouponByIds(joinCouponIds)
	}
	redPacket["couponMaps"] = couponMaps
	receiveRedPacket, err := GetReceiveRedPacket(ws.Unionid, redPacketId)
	if err != nil {
		return utils.ErrorNull(c, "获取用户的红包领取记录失败")
	}
	redPacket["receiveRedPacket"] = receiveRedPacket
	return utils.SuccessNullMsg(c, redPacket)
}

func GetRedPacket(id int64) map[string]interface{} {
	m, err := global.DB.Query("SELECT * FROM account_send_red_packet WHERE status=? AND id=? LIMIT 1", enum.NORMAL, id)
	if err != nil {
		global.Log.Error("GetRedPacket sql ERROR：", err.Error())
		return nil
	}
	if len(m) != 1 {
		return nil
	}
	return m[0]
}

/**
* 发送红包
 */
func SendRedPacket(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	sendName := strings.TrimSpace(c.FormValue("sendName"))
	sendPhoto := strings.TrimSpace(c.FormValue("sendPhoto"))
	moneyStr := strings.TrimSpace(c.FormValue("money"))
	randomMoney := strings.TrimSpace(c.FormValue("randomMoney"))
	blessings := strings.TrimSpace(c.FormValue("blessings"))
	remarks := strings.TrimSpace(c.FormValue("remarks"))
	joinCouponIds := strings.TrimSpace(c.FormValue("joinCouponIds"))
	//payType := strings.TrimSpace(c.FormValue("payType"))
	if sendName == "" {
		return utils.ErrorNull(c, "请填写发送红包人的名称")
	}
	if sendPhoto == "" {
		return utils.ErrorNull(c, "请上传发送红包人的头像")
	}
	if moneyStr == "" || !utils.IsValidNumber(moneyStr) {
		return utils.ErrorNull(c, "请正确填写红包金额")
	}
	money := convert.MustFloat64(moneyStr)
	if money < 10 {
		return utils.ErrorNull(c, "红包金额最低10元")
	}
	if randomMoney == "" || !utils.IsValidNumber(randomMoney) {
		return utils.ErrorNull(c, "请正确选择随机金额范围")
	}
	if convert.MustFloat64(randomMoney) > convert.MustFloat64(moneyStr) {
		return utils.ErrorNull(c, "随机金额最大值不能超过红包的总金额")
	}
	if blessings != "" && len(blessings) > 100 {
		return utils.ErrorNull(c, "祝福语字数过长，最长不超过100字")
	}
	if remarks != "" && len(remarks) > 100 {
		return utils.ErrorNull(c, "备注字数过长，最长不超过500字")
	}
	if joinCouponIds != "" {
		couponIds := strings.Split(joinCouponIds, ",")
		if len(couponIds) > 0 {
			for i := 0; i < len(couponIds); i++ {
				if !utils.IsValidNumber(couponIds[i]) {
					return utils.ErrorNull(c, "选择代金券错误")
				}
			}
			//判断代金券有效性
			couponMaps := GetCouponByIds(joinCouponIds)
			if couponMaps == nil || len(couponIds) != len(couponMaps) {
				return utils.ErrorNull(c, "选择代金券错误或代金券已删除")
			}
		}
	}
	id := utils.ID()
	currentTime := utils.CurrentTime()
	ip := c.RealIP()
	_, err = global.DB.InsertMap("account_send_red_packet", map[string]interface{}{
		"id":              id,
		"account_id":      acc.ID,
		"send_name":       sendName,
		"send_photo":      sendPhoto,
		"money":           money,
		"fee_money":       money * 0.02,
		"random_money":    randomMoney,
		"blessings":       blessings,
		"remarks":         remarks,
		"pay_state":       enum.PAY_STATUS_NOT,
		"ct_time":         currentTime,
		"ut_time":         currentTime,
		"join_coupon_ids": joinCouponIds,
		"status":          enum.NORMAL,
		"ip":              ip,
	})
	if err != nil {
		global.Log.Error("global.DB.InsertMap account_send_red_packet sql error:%s", err.Error())
		return utils.ErrorNull(c, "创建红包失败")
	} else {
		return utils.Success(c, "创建红包成功", id)
	}
}

/*
*发送红包记录
 */
func GetSendRedPacketList(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	t := c.FormValue("t")
	//已支付，已完成
	where := fmt.Sprintf("pay_state='%s' AND money-send_money>=1", enum.PAY_STATUS_END)
	switch t {
	case "finish":
		//已领完
		where = fmt.Sprintf("pay_state='%s' AND money-send_money<1", enum.PAY_STATUS_END)
		break
	case "notPay":
		//未支付
		where = fmt.Sprintf("pay_state='%s'", enum.PAY_STATUS_NOT)
		break
	}
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "account_send_red_packet",
		Where:     "account_id=? AND status=? AND " + where,
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "ct_time DESC",
	}, acc.ID, enum.NORMAL)
	if err != nil {
		global.Log.Error("GetSendRedPacketList sql error:%s", err.Error())
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}
