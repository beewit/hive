/*
发送红包记录
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
* [迁移至pay项目]发送红包
 */
func SendRedPacket(c echo.Context) error {
	acc, err := GetAccount(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	sendName := strings.TrimSpace(c.FormValue("sendName"))
	sendPhoto := strings.TrimSpace(c.FormValue("sendPhoto"))
	money := strings.TrimSpace(c.FormValue("money"))
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
	if money == "" || !utils.IsValidNumber(money) {
		return utils.ErrorNull(c, "请正确填写红包金额")
	}
	if randomMoney == "" || !utils.IsValidNumber(randomMoney) {
		return utils.ErrorNull(c, "请正确选择随机金额范围")
	}
	if convert.MustFloat64(randomMoney) > convert.MustFloat64(money) {
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
		return utils.ErrorNull(c, "创建红包失败")
	} else {
		return utils.Success(c, "创建红包成功", id)
	}
}
