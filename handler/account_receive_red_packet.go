/*
领取红包记录
author：zxb
2018-01-13
*/

package handler

import (
	"fmt"
	"github.com/beewit/beekit/utils"
	"github.com/beewit/beekit/utils/convert"
	"github.com/beewit/beekit/utils/enum"
	"github.com/beewit/beekit/utils/uhttp"
	"github.com/beewit/hive/global"
	"github.com/labstack/echo"
	"strings"
	"time"
)

func ReceiveRedPacket(c echo.Context) error {
	ws, err := GetMiniAppSession(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := c.FormValue("id")
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	redPacketId := convert.MustInt64(id)
	redPacket := GetRedPacket(redPacketId)
	if redPacket == nil {
		return utils.ErrorNull(c, "该红包不存在")
	}
	sendRedPacketAcc := GetAccountById(convert.MustInt64(redPacket["account_id"]))
	if sendRedPacketAcc == nil {
		return utils.ErrorNull(c, "发红包的账号已被锁定无法领取")
	}
	if redPacket["status"] != enum.NORMAL {
		return utils.ErrorNull(c, "该红包已失效")
	}
	if redPacket["pay_state"] != enum.PAY_STATUS_END {
		return utils.ErrorNull(c, "红包无效")
	}
	receiveRedPacket, err := GetReceiveRedPacket(ws.Unionid, redPacketId)
	if err != nil {
		return utils.ErrorNull(c, "获取用户的红包领取记录失败")
	}
	if receiveRedPacket != nil {
		createQrCode := false
		qrCodeTime := convert.ToString(receiveRedPacket["qrcode_time"])
		if convert.ToString(receiveRedPacket["qrcode"]) != "" && qrCodeTime != "" {
			//判断是否过期
			qrcodeTime, err := time.Parse("2006-01-02 15:04:05", qrCodeTime)
			if err != nil {
				global.Log.Error("领红包二维码到期时间异常：%s", err.Error())
				return utils.ErrorNull(c, "领红包二维码到期时间异常：")
			}
			if !qrcodeTime.After(time.Now()) {
				//已过期
				createQrCode = true
			}
		} else {
			//重新生成
			createQrCode = true
		}
		if createQrCode {
			qrCodePath := UpdateRedPacketQrCode(convert.MustInt64(receiveRedPacket["id"]))
			if qrCodePath != "" {
				receiveRedPacket["qrcode"] = qrCodePath
			}
		}
		return utils.Success(c, "已领取过次红包，不再添加记录", receiveRedPacket)
	}
	money := convert.MustFloat64(redPacket["money"])
	sendMoney := convert.MustFloat64(redPacket["send_money"])
	randomMoney := convert.MustFloat64(redPacket["random_money"])
	if money < 1 {
		return utils.ErrorNull(c, "该红包已领完")
	}
	if money-sendMoney < randomMoney {
		return utils.ErrorNull(c, "该红包已领完")
	}

	acc := GetAccountByUnionId(ws.Unionid, enum.WECHAT)
	var accId interface{}
	if acc != nil {
		accId = acc["id"]
	} else {
		accId = nil
	}
	receiveRedPacketId := utils.ID()
	receiveRedPacket = map[string]interface{}{
		"id":                         receiveRedPacketId,
		"wx_union_id":                ws.Unionid,
		"account_id":                 accId,
		"account_send_red_packet_id": redPacket["id"],
		"ct_time":                    utils.CurrentTime(),
		"ip":                         c.RealIP(),
		"status":                     enum.RED_PACKET_STATUS_NOT,
	}
	_, err = global.DB.InsertMap("account_receive_red_packet", receiveRedPacket)
	if err != nil {
		return utils.ErrorNull(c, "领取红包失败")
	}
	qrCodePath := UpdateRedPacketQrCode(receiveRedPacketId)
	if qrCodePath != "" {
		receiveRedPacket["qrcode"] = qrCodePath
	}
	return utils.Success(c, "领取红包成功", receiveRedPacket)
}

func UpdateRedPacketQrCode(receiveRedPacketId int64) string {
	global.Log.Info("生成领取红包二维码！")
	var qrCodePath string
	body, err := uhttp.Cmd(uhttp.Request{
		Method: "POST",
		URL:    fmt.Sprintf("http://m.9ee3.com/account/create/temporary/qrcode?objId=%v&objType=%s", receiveRedPacketId, enum.QRCODE_RED_PACKET),
	})
	if err != nil {
		global.Log.Error("获取领取红包临时二维码失败，%v", err.Error())
	} else {
		global.Log.Info(string(body))
		resultParam := utils.ToResultParam(body)
		if resultParam.Ret == utils.SUCCESS_CODE {
			data, err := convert.Obj2Map(resultParam.Data)
			if err != nil {
				global.Log.Error("获取领取红包临时二维码失败，转换数据失败：%v", err.Error())
				return ""
			} else {
				//保存
				qrCodePath = convert.ToString(data["path"])
			}
		} else {
			global.Log.Error("获取领取红包临时二维码失败，%v", resultParam.Msg)
		}
	}
	if qrCodePath != "" {
		x, err := global.DB.Update("UPDATE account_receive_red_packet SET qrcode=?,qrcode_time=? WHERE id=?", qrCodePath,
			utils.FormatTime(time.Now().Add(-time.Hour)), receiveRedPacketId)
		if err != nil {
			global.Log.Error(err.Error())
			return ""
		}
		if x <= 0 {
			global.Log.Error(fmt.Sprintf("%v修改红包二维码失败", receiveRedPacketId))
		}
	}
	global.Log.Info("【结果】生成领取红包二维码：%s", qrCodePath)
	return qrCodePath
}

//获取红包领取记录
func GetReceiveRedPacket(uinonId string, redPacketId int64) (map[string]interface{}, error) {
	sql := "SELECT * FROM account_receive_red_packet WHERE wx_union_id=? AND account_send_red_packet_id=? LIMIT 1"
	rows, err := global.DB.Query(sql, uinonId, redPacketId)
	if err != nil {
		global.Log.Error("IsReceiveRedPacket sql error:", err.Error())
		return nil, err
	}
	if len(rows) != 1 {
		return nil, nil
	}
	return rows[0], nil
}

/**
*领取红包记录
 */
func GetReceiveRedPacketList(c echo.Context) error {
	ws, err := GetMiniAppSession(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	pageIndex := utils.GetPageIndex(c.FormValue("pageIndex"))
	pageSize := utils.GetPageSize(c.FormValue("pageSize"))
	page, err := global.DB.QueryPage(&utils.PageTable{
		Fields:    "*",
		Table:     "v_account_receive_red_packet",
		Where:     "wx_union_id=?",
		PageIndex: pageIndex,
		PageSize:  pageSize,
		Order:     "receiveTime DESC",
	}, ws.Unionid)
	if err != nil {
		global.Log.Error("GetReceiveRedPacketList v_account_receive_red_packet sql error:%s", err.Error())
		return utils.Error(c, "数据异常，"+err.Error(), nil)
	}
	if page == nil {
		return utils.NullData(c)
	}
	return utils.Success(c, "获取数据成功", page)
}

//领取红包和优惠券的记录
func GetReceiveRedPacketAndCouponList(c echo.Context) error {
	_, err := GetMiniAppSession(c)
	if err != nil {
		return utils.AuthFailNull(c)
	}
	id := strings.TrimSpace(c.FormValue("id"))
	if id == "" || !utils.IsValidNumber(id) {
		return utils.ErrorNull(c, "id格式错误")
	}
	redPacket := GetRedPacket(convert.MustInt64(id))
	if redPacket == nil {
		return utils.ErrorNull(c, "红包不存在或已过期")
	}
	sql := "SELECT ar.money,ar.ct_time as receiveTime,wa.* FROM account_receive_red_packet ar LEFT JOIN wx_account wa ON ar.wx_union_id =wa.union_id WHERE account_send_red_packet_id=? AND status<>?"
	redPacketList, err := global.DB.Query(sql, id, enum.RED_PACKET_STATUS_NOT)
	if err != nil {
		global.Log.Error("GetReceiveRedPacketAndCouponList account_receive_red_packet sql error:%s", err.Error())
		return utils.Error(c, "获取领取红包数据失败", nil)
	}
	couponList := []map[string]interface{}{}
	joinCouponIds := convert.ToString(redPacket["join_coupon_ids"])
	if joinCouponIds != "" {
		sql = fmt.Sprintf("SELECT ac.money,ar.ct_time as receiveTime,wa.* FROM account_receive_coupon ar LEFT JOIN wx_account wa ON ar.wx_union_id =wa.union_id "+
			"LEFT JOIN account_coupon ac ON ac.id=ar.account_coupon_id WHERE account_coupon_id in(%s)", joinCouponIds)
		couponList, err = global.DB.Query(sql)
		if err != nil {
			global.Log.Error("GetReceiveRedPacketAndCouponList account_receive_coupon sql error:%s", err.Error())
			return utils.Error(c, "获取领取代金券数据失败", nil)
		}
	}
	return utils.SuccessNullMsg(c, map[string]interface{}{
		"redPacket":     redPacket,
		"redPacketList": redPacketList,
		"couponList":    couponList,
	})
}
